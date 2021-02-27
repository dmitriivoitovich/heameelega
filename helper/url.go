package helper

import (
	"net/url"
	"strconv"

	"github.com/dmitriivoitovich/heameelega/config"
	"github.com/google/uuid"
)

func PageURLHome() string {
	scheme := "http://"
	if config.AppTLS().Enabled {
		scheme = "https://"
	}

	return scheme + config.AppHost() + "/"
}

func PageURLLogin() string {
	return PageURLHome() + "login"
}

func PageURLRegister() string {
	return PageURLHome() + "register"
}

func PageURLLogout() string {
	return PageURLHome() + "logout"
}

func PageURLUserSettings() string {
	return PageURLHome() + "settings"
}

func PageURLDashboard() string {
	return PageURLHome() + "dashboard"
}

func PageURLCustomers() string {
	return PageURLHome() + "customers"
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

	return PageURLHome() + "customers?" + params.Encode()
}

func PageURLNewCustomer() string {
	return PageURLHome() + "customers/new"
}

func PageURLViewCustomer(id uuid.UUID) string {
	return PageURLHome() + "customers/" + id.String()
}

func PageURLEditCustomer(id uuid.UUID) string {
	return PageURLViewCustomer(id) + "/edit"
}

func PageURLDeleteCustomer(id uuid.UUID) string {
	return PageURLViewCustomer(id) + "/delete"
}
