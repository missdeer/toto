package alipay

import (
	"crypto/md5"
	"fmt"
	"github.com/missdeer/toto/routers/base"
	"github.com/missdeer/toto/setting"
	"net/url"
	"sort"
	"strings"
)

// HomeRouter serves home page.
type AlipayRouter struct {
	base.BaseRouter
}

func (this *AlipayRouter) Pay() {
	finalUrl := `https://mapi.alipay.com/gateway.do?`

	sign_type := "MD5"
	partner := setting.AlipayPartnerId
	sign_key := setting.AlipaySignKey
	seller_email := setting.AlipaySellerEmail
	seller_id := setting.AlipaySellerId
	service := "trade_create_by_buyer"
	_input_charset := "utf-8"
	out_trade_no := this.GetString("out_trade_no")
	subject := this.GetString("subject")
	price := this.GetString("price")
	quantity := this.GetString("quantity")
	payment_type := "1"
	logistics_type := this.GetString("logistics_type")
	logistics_fee := this.GetString("logistics_fee")
	logistics_payment := this.GetString("logistics_payment")
	//body := this.GetString("body")
	//notify_url := `http://httpapi.sinaapp.com/alipay.php?action=notify`

	var ss sort.StringSlice
	//ss = append(ss, fmt.Sprintf("body=%s", body))
	ss = append(ss, fmt.Sprintf("partner=%s", partner))
	//ss = append(ss, fmt.Sprintf("sign_key=%s", sign_key))
	//ss = append(ss, fmt.Sprintf("notify_url=%s", notify_url))
	if len(seller_email) > 0 {
		ss = append(ss, fmt.Sprintf("seller_email=%s", seller_email))
	}
	if len(seller_id) > 0 {
		ss = append(ss, fmt.Sprintf("seller_id=%s", seller_id))
	}
	ss = append(ss, fmt.Sprintf("service=%s", service))
	ss = append(ss, fmt.Sprintf("_input_charset=%s", _input_charset))
	ss = append(ss, fmt.Sprintf("out_trade_no=%s", out_trade_no))
	ss = append(ss, fmt.Sprintf("subject=%s", subject))
	ss = append(ss, fmt.Sprintf("price=%s", price))
	ss = append(ss, fmt.Sprintf("quantity=%s", quantity))
	ss = append(ss, fmt.Sprintf("logistics_type=%s", logistics_type))
	ss = append(ss, fmt.Sprintf("logistics_fee=%s", logistics_fee))
	ss = append(ss, fmt.Sprintf("logistics_payment=%s", logistics_payment))
	ss = append(ss, fmt.Sprintf("payment_type=%s", payment_type))
	ss.Sort()
	parameterString := strings.Join(ss, "&")
	fmt.Println("parameter string: ", parameterString)
	readyString := fmt.Sprintf("%s%s", parameterString, sign_key)
	fmt.Println("ready string: ", readyString)
	sum := md5.Sum([]byte(readyString))
	sign := fmt.Sprintf("%032x", sum)
	fmt.Println("sign result: ", sign)
	fmt.Println("a md5 sum: ", fmt.Sprintf("%032x", md5.Sum([]byte("a"))))

	var ff sort.StringSlice
	//ff = append(ff, fmt.Sprintf("body=%s", url.QueryEscape(body)))
	ff = append(ff, fmt.Sprintf("partner=%s", url.QueryEscape(partner)))
	//ff = append(ff, fmt.Sprintf("sign_key=%s", url.QueryEscape(sign_key)))
	if len(seller_email) > 0 {
		ff = append(ff, fmt.Sprintf("seller_email=%s", url.QueryEscape(seller_email)))
	}
	if len(seller_id) > 0 {
		ff = append(ff, fmt.Sprintf("seller_id=%s", url.QueryEscape(seller_id)))
	}
	ff = append(ff, fmt.Sprintf("service=%s", url.QueryEscape(service)))
	ff = append(ff, fmt.Sprintf("_input_charset=%s", url.QueryEscape(_input_charset)))
	ff = append(ff, fmt.Sprintf("out_trade_no=%s", url.QueryEscape(out_trade_no)))
	ff = append(ff, fmt.Sprintf("subject=%s", url.QueryEscape(subject)))
	ff = append(ff, fmt.Sprintf("price=%s", url.QueryEscape(price)))
	ff = append(ff, fmt.Sprintf("quantity=%s", url.QueryEscape(quantity)))
	ff = append(ff, fmt.Sprintf("logistics_type=%s", url.QueryEscape(logistics_type)))
	ff = append(ff, fmt.Sprintf("logistics_fee=%s", url.QueryEscape(logistics_fee)))
	ff = append(ff, fmt.Sprintf("logistics_payment=%s", url.QueryEscape(logistics_payment)))
	ff = append(ff, fmt.Sprintf("payment_type=%s", url.QueryEscape(payment_type)))
	//ff = append(ff, fmt.Sprintf("notify_url=%s", url.QueryEscape(notify_url)))
	ff.Sort()
	ff = append(ff, fmt.Sprintf("sign=%s", url.QueryEscape(sign)))
	ff = append(ff, fmt.Sprintf("sign_type=%s", url.QueryEscape(sign_type)))
	parameterString = strings.Join(ff, "&")

	finalUrl = finalUrl + parameterString
	fmt.Println("final url: ", finalUrl)

	this.Redirect(finalUrl, 301)
}
