package helper

import (
	"net/url"
	"strconv"

	"github.com/dmitriivoitovich/heameelega/config"
	"github.com/google/uuid"
)

func PageURLHome() string {
	return config.AppHost() + "/"
}

func PageURLLogout() string {
	return config.AppHost() + "/logout"
}

func PageURLDashboard() string {
	return config.AppHost() + "/dashboard"
}

func PageURLCustomers() string {
	return config.AppHost() + "/customers"
}

func PageURLSearchCustomers(filter, order, direction string, reverseDirection bool, page uint32) string {
	if reverseDirection {
		if direction == "asc" {
			direction = "desc"
		} else {
			direction = "asc"
		}
	}

	params := url.Values{
		"o": {order},
		"d": {direction},
		"p": {strconv.Itoa(int(page))},
	}

	if filter != "" {
		params.Add("s", filter)
	}

	return config.AppHost() + "/customers?" + params.Encode()
}

func PageURLNewCustomer() string {
	return config.AppHost() + "/customers/new"
}

func PageURLViewCustomer(id uuid.UUID) string {
	return config.AppHost() + "/customers/" + id.String()
}

func PageURLEditCustomer(id uuid.UUID) string {
	return config.AppHost() + "/customers/" + id.String() + "/edit"
}
