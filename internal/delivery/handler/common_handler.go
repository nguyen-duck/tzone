package handler

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/LuuDinhTheTai/tzone/util/response"
	"github.com/LuuDinhTheTai/tzone/web/page"
	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

func (h *CommonHandler) IndexHandler(ctx *gin.Context) {
	response.HTML(ctx, page.HomePage())
}

func (h *CommonHandler) SignupHandler(ctx *gin.Context) {
	response.HTML(ctx, page.SignupPage())
}

func (h *CommonHandler) AllBrandHandler(ctx *gin.Context) {
	response.HTML(ctx, page.AllBrandPage())
}

func (h *CommonHandler) BrandHandler(ctx *gin.Context) {
	brand := ctx.Query("brand")
	if brand == "" {
		brand = "acer"
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageIdx, err := strconv.Atoi(pageStr)
	if err != nil || pageIdx < 1 {
		pageIdx = 1
	}

	allDevices := GenerateMockDevices(brand)
	total := len(allDevices)
	perPage := 50
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	if pageIdx > totalPages && totalPages > 0 {
		pageIdx = totalPages
	}

	start := (pageIdx - 1) * perPage
	end := start + perPage
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedDevices := allDevices[start:end]

	brandCapitalized := "Acer"
	if len(brand) > 0 {
		brandCapitalized = strings.ToUpper(brand[:1]) + strings.ToLower(brand[1:])
	}
	response.HTML(ctx, page.BrandPage(brandCapitalized, pagedDevices, pageIdx, totalPages))
}

func GenerateMockDevices(brand string) []page.Device {
	var seed = []page.Device{
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-Iconia-v12-.jpg", Name: "Iconia V12"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-iconia-v11-.jpg", Name: "Iconia V11"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-acerone-liquid-s262f5.jpg", Name: "Acerone Liquid S262F5"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-iconia-a14-.jpg", Name: "Iconia A14"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-iconia-x14-.jpg", Name: "Iconia X14"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-iconia-x12.jpg", Name: "Iconia X12"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-iconia-tab-p11.jpg", Name: "Iconia Tab P11 (2025)"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-iconia-tab-p10.jpg", Name: "Iconia Tab P10 (2025)"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-super-zx-pro-1.jpg", Name: "Super ZX Pro"},
		{ImgSrc: "https://fdn2.gsmarena.com/vv/bigpic/acer-super-zx.jpg", Name: "Super ZX"},
	}

	count := 0
	brandLower := strings.ToLower(brand)
	if brandLower == "acer" {
		count = 113
	} else if brandLower == "amoi" {
		count = 47
	} else {
		count = 0
	}

	var devices []page.Device
	for i := 0; i < count; i++ {
		seedDev := seed[i%len(seed)]
		name := seedDev.Name
		if brandLower != "acer" {
			name = fmt.Sprintf("%s Device %d", strings.ToUpper(brandLower[:1])+brandLower[1:], i+1)
		} else {
			if i >= len(seed) {
				name = fmt.Sprintf("%s (V%d)", name, (i/len(seed))+1)
			}
		}
		devices = append(devices, page.Device{
			Name:   name,
			ImgSrc: seedDev.ImgSrc,
		})
	}
	return devices
}
