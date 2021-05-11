package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateTransaction(t *testing.T) {

	samples := []struct {
		statusCode    int
		productID     string
		affiliateCode string
		error         string
		message       string
	}{
		{
			statusCode:    200,
			productID:     seededProduct.ID,
			affiliateCode: seededAffiliate.AffiliateCode,
			error:         "",
			message:       "purchase successful",
		},
		{
			statusCode:    400,
			productID:     "invalid_product_id",
			affiliateCode: seededAffiliate.AffiliateCode,
			error:         "invalid product provided",
			message:       "",
		},
		{
			statusCode:    400,
			productID:     seededProduct.ID,
			affiliateCode: "invalid_affiliate_id",
			error:         "invalid affiliate provided",
			message:       "",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/purchase-product/%s/%s", server.URL, v.productID, v.affiliateCode), nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		responseMap := make(map[string]string)
		err = json.Unmarshal(body, &responseMap)
		assert.NoError(t, err)

		assert.Equal(t, resp.StatusCode, v.statusCode)
		if resp.StatusCode == 200 {
			assert.Equal(t, responseMap["message"], v.message)
		} else if resp.StatusCode > 299 {
			assert.Equal(t, responseMap["error"], v.error)
		}
	}
}
