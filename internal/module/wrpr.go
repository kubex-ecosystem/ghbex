package module

import (
	"os"
	"strings"
)

func RegX() *Ghbex {
	var printBannerV = os.Getenv("GHBEX_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &Ghbex{
		PrintBanner: strings.ToLower(printBannerV) == "true",
	}
}
