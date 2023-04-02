package controllers

import (
	"context"
	"github.com/FelixWuu/chatgpt_self_use/config"
	"github.com/FelixWuu/chatgpt_self_use/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	goGpt "github.com/sashabaranov/go-gpt3"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	socks5Position = 10
)

type Response struct{}
type JsonController struct{}

func (j *JsonController) GptRspJson(ctx *gin.Context, code int, errorMsg string, data interface{}) {
	ctx.JSON(code, gin.H{
		"code":     code,
		"errorMsg": errorMsg,
		"data":     data,
	})
	ctx.Abort()
}

type ResponseController struct {
	JsonController
}

func NewResponseController() *ResponseController {
	return &ResponseController{}
}

// Index 首页
func (r *ResponseController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main Website",
	})
}

// Response 响应体，ChatGPT 用于回复用户的消息
func (r *ResponseController) Response(ctx *gin.Context) {
	// 1. 封装发给 gpt 的请求，不允许发送空请求
	var request goGpt.ChatCompletionRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		r.GptRspJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logger.Info(request)
	if len(request.Messages) == 0 {
		r.GptRspJson(ctx, http.StatusBadRequest, "missing request information", nil)
		return
	}

	// 2. 设置代理，防墙
	cfg := config.Inst()
	gptCfg := goGpt.DefaultConfig(cfg.Api.ApiKey)
	if err := setProxyService(cfg.Address.Proxy, &gptCfg); err != nil {
		r.GptRspJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 3. 允许自定义 BaseURL
	apiURL := cfg.Api.ApiUrl
	if apiURL != "" {
		gptCfg.BaseURL = apiURL
	}

	// 4. 发送请求并接收 ChatGPT 的消息
	client := r.createClient(&gptCfg, request, cfg.Bot.Personalization)
	r.requestChatCompletion(ctx, request, client, cfg)
}

// setProxyService 设置代理，防止被墙
func setProxyService(proxyAddr string, gptCfg *goGpt.ClientConfig) error {
	if proxyAddr != "" {
		transport := &http.Transport{}
		if strings.HasPrefix(proxyAddr, "socks5h://") {
			// 设置一个 dailContext 对象的代理服务器
			dialCtx, err := newDialContext(proxyAddr[socks5Position:])
			if err != nil {
				return errors.Wrap(err, "create dialContext failed")
			}

			transport.DialContext = dialCtx
		} else {
			// 设置一个 http transport 对象的代理服务器
			proxyUrl, err := url.Parse(proxyAddr)
			if err != nil {
				return errors.Wrap(err, "create http transport failed")
			}

			transport.Proxy = http.ProxyURL(proxyUrl)
		}

		gptCfg.HTTPClient = &http.Client{
			Transport: transport,
		}
	}

	return nil
}

type dialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

// newDialContext 新建一个 dail context
func newDialContext(socks5 string) (dialContextFunc, error) {
	baseDialer := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}

	if socks5 != "" {
		// split socks5 proxy string [username:password@]host:port
		var auth *proxy.Auth = nil

		if strings.Contains(socks5, "@") {
			proxyInfo := strings.SplitN(socks5, "@", 2)
			proxyUser := strings.Split(proxyInfo[0], ":")
			if len(proxyUser) == 2 {
				auth = &proxy.Auth{
					User:     proxyUser[0],
					Password: proxyUser[1],
				}
			}
			socks5 = proxyInfo[1]
		}

		dialSocksProxy, err := proxy.SOCKS5("tcp", socks5, auth, baseDialer)
		if err != nil {
			return nil, err
		}

		contextDialer, ok := dialSocksProxy.(proxy.ContextDialer)
		if !ok {
			return nil, err
		}

		return contextDialer.DialContext, nil
	} else {
		return baseDialer.DialContext, nil
	}
}

// createClient 创建 ChatGPT 客户端
func (r *ResponseController) createClient(
	gptCfg *goGpt.ClientConfig,
	request goGpt.ChatCompletionRequest,
	personalization string,
) *goGpt.Client {
	client := goGpt.NewClientWithConfig(*gptCfg)
	if request.Messages[0].Role != "system" {
		newMessage := append([]goGpt.ChatCompletionMessage{
			{Role: "system", Content: personalization},
		}, request.Messages...)
		request.Messages = newMessage
		logger.Info(request.Messages)
	}
	return client
}

func (r *ResponseController) requestChatCompletion(
	ctx *gin.Context,
	request goGpt.ChatCompletionRequest,
	client *goGpt.Client,
	cfg *config.Config,
) {
	model := cfg.Bot.Model

	if model == goGpt.GPT3Dot5Turbo || model == goGpt.GPT3Dot5Turbo0301 {
		request.Model = model
		rsp, err := client.CreateChatCompletion(ctx, request)
		if err != nil {
			r.GptRspJson(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		r.GptRspJson(ctx, http.StatusOK, "", gin.H{
			"reply":   rsp.Choices[0].Message.Content,
			"message": append(request.Messages, rsp.Choices[0].Message),
		})
	} else {
		prompt := ""
		for _, msg := range request.Messages {
			prompt += msg.Content + "/n"
		}

		prompt = strings.Trim(prompt, "/n")
		logger.Info("request prompt: %s", prompt)

		cmplReq := goGpt.CompletionRequest{
			Model:            model,
			MaxTokens:        cfg.Bot.MaxTokens,
			TopP:             cfg.Bot.TopP,
			FrequencyPenalty: cfg.Bot.FrequencyPenalty,
			PresencePenalty:  cfg.Bot.PresencePenalty,
			Prompt:           prompt,
		}

		rsp, err := client.CreateCompletion(ctx, cmplReq)
		if err != nil {
			r.GptRspJson(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		r.GptRspJson(ctx, http.StatusOK, "", gin.H{
			"reply": rsp.Choices[0].Text,
			"message": append(request.Messages, goGpt.ChatCompletionMessage{
				Role:    "assistant",
				Content: rsp.Choices[0].Text,
			}),
		})
	}
}
