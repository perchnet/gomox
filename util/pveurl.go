package util

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const ApiUrlSuffix = "/api2/json"

/*
GetPveUrl returns either the URL as specified by the `pveurl` arg,
or builds a URL from the `scheme`, `pvehost`, and `pveport` args.
*/
func GetPveUrl(c *cli.Context) string {
	var ret string
	pveUrl := c.String("pveurl")
	switch pveUrl {
	case "":
		ret = fmt.Sprint(
			c.String("scheme"),
			"://",
			c.String("pvehost"),
			":",
			c.String("pveport"),
			ApiUrlSuffix,
		)
	default:
		ret = pveUrl
	}
	return ret
}
