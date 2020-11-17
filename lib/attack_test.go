package happa

import (
	"net/http"
	"testing"
	"time"
)

func TestCalcWaitTime(t *testing.T) {
	c := http.Client{}
	url := ""
	method := ""
	var headers []string
	duration := 0
	workerNumber := 0
	for _, test := range []struct {
		elapsed time.Duration
		rate int
		rc int
		expected time.Duration
	} {
		{1, -1000000, 0, 0},
		{1, -1, 0, 0},
		{1, 0, 0, 0},
		{1, 1, 0, 0},
		{1, 1, 1, 1000000000},
		{1, 1000, 999, 0},
		{1, 1000, 1001, 1000000},
		{1, 1000000, 999999, 0},
		{1, 1000000, 1000001, 1000},
		{1, 1000000000, 999999999, 0},
		{1, 1000000000, 1000000001, 1},
		{1, 1000000001, 1000000000, 0},
		{1, 1000000001, 1000000002, 1},
	} {
		a := NewAttacker(&c, &url, &method, headers, &duration, &test.rate, &workerNumber)
		actual := a.calcWaitTime(&test.rc, &test.elapsed)
		if actual != test.expected {
			t.Errorf("[TestCalcWaitTime] actual = %d, expected = %d", actual, test.expected)
		}
	}
}

func TestParseHeader(t *testing.T) {
	url := ""
	method := ""
	duration := 0
	rate := 0
	workerNumber := 0
	for _, test := range []struct {
		headers Headers
		expected map[string][]string
	} {
		{
			[]string{
				"authority: www.yahoo.co.jp",
				"cache-control: max-age=0",
				"upgrade-insecure-requests: 1",
				"user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
				"accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
				"sec-fetch-site: none",
				"sec-fetch-mode: navigate",
				"sec-fetch-user: ?1",
				"sec-fetch-dest: document",
				"accept-language: ja,en-US;q=0.9,en;q=0.8",
				"cookie: F=a=AjGEC_IMvScXQ6jRpdCiF1Zhsv9NxTEB3jReaRXcj4fC9mnYtTgO5u5ZLB8FAmstQIK8ub4-&b=pZ98; B=1l742shfnu3kk&b=4&d=f.PR2GlpYF0ffQzA9.sN9PkjdO1k6F_IfPuvDYwl&s=0k&i=16bahJ2QKNr6oOzXHAKD; XB=1l742shfnu3kk&b=4&d=f.PR2GlpYF0ffQzA9.sN9PkjdO1k6F_IfPuvDYwl&s=0k&i=16bahJ2QKNr6oOzXHAKD; Y=v=1&n=bj8bbcaff7t61&l=c0c0e_o/o&p=m2kvvjp052000000&ig=01767&r=115&lg=ja-JP&intl=jp; T=z=BA4gfBBoGqfBIXJK3.cWRIkMDc3MwYxNjFOMDUwTzA-&sk=DAAwvb7OXbcH80&ks=EAAHuBJ4plFljvAGaxRTAEuEg--~F&kt=EAAD0Pb3pbZ.h4xxZ8GbOgmbw--~E&ku=FAAMEYCIQCrcvvQ2x.iLqbg_A1F25vdOsdlKn7vNCl4CPRA.GjH7QIhAIWRpIOGK2H8R1fW7g_LLf8mE9ekXSPP8FIvBByI7_Y3~B&d=dGlwATJGeDBEQwFhAVlBRQFnATRTUVJLS1U0WERWRVhWU0ZENFBVNFNCWE5BAXNsAU56QXdOQUUyTVRZNU56STNPRGMtAXNjAXd3dwF6egFCQTRnZkJBMko-; SSL=v=1&s=H8a9563jCll2MrcDpijrxS8QQ2tVAOJnetfb4_UMEiEdQThv5tHTfRUHBhp18zgi.wAgwzbYeLgKffvW9yJXlw--&kv=0; _n=DPFnxL-4PtCUH4M4G1idsx8MuEQMPWcLNUO8kFSRYa9MTwBAZjRLGMs4GnGMVjk4IEKcOFHNMdXMF-goe3DvrEbYbvnErNZIzgZtBdAP-FIPPE1O8QeolMcT6Al37amXltSx7tR7011AcqLTsPE7yDPkgcrGKPwBubo6vjDJp2GYHPfkDaWhTRyWowvzKG7b3GAs38apZuiUCk4chVfNSJZoB_IeuJ4rohz3m6aYWSm2NWiKe-0xwaDJebgmhI7e16ab3lu8TLqR2wzqTq0j_jIXMEl49l0bL8fBeH87gYvXcsjvape1d3l3QeVRb1SwgfVCLj3QJKAvowjUydSGYBUemL2j1Rrb6gm_Cbf4Vb8kWDJVcE1abCPQG-N5uPobJA1DEl623ihTOfSePvJIAibHDs2MUfAklkAVAqO530Ha66SBl1lZZRJIxy1IFigSWTLknJHEuvj981Qehwv8E63r6KyinP89GpN_fyMh1h8YVki6NQhvp0cbN8jfc9yPHleV-9t8IpJHDZzUKS3Db6NNDDZdo1LhM63xrCxbDMpdeNFtahRkA2KfwmSMuZMu0CO2x-pwHAhgLoUUFP_Yik4hmY5gGnQZasrsTk_2ZJ8U4fLhyGcgnxhfZtvPMz5qIYKwHDMOm6vf5C_ql9-HJjR_zedsb5IkIVf2zE1ZVBKUhjLwrmJNof93Ag2dH3h4CFRTaGlfoxNlWJe0S5yuNnp6YUxqMxAIRf6HwXCnKdKGD07zblPPqvte9wALQpOmf6BwmpaFpqb-VP1-qfj3WVv2gX7HAa_G75mFVush-mtP45kSO4TCBG9j9MutYOHCkjYFw3V3HCY9UaLeUY56mHrK63ySJPuMXrX7_lKJemYcbvId4TsEo_kuufGWM-RueWM8EHvpvS8UooQveOwRN15hEJL-Zdi_8JHZOlCegHovMbJWLY3avtVJbmip_oV_p6LNlSWkeM8E5RyBQuPwZskXjSX21Sv2gOXSCVmyNI-RYWlIEnd2tIfMaU5gGyUnMBJTIqDIP-nsr9AjkEslJAJPUCX_mMtRM8_EGhoX1OS8m_jqDNcJBzR1dE1KFN8du8Xb1cEfPHRy8V2mHHC8FavMlgS9esUgzELkU6vQYnVs_-YGJsP1osFFx5ebk_c8.2; XC=JnDKyz8lyzvMSPDcpm99NAG3Cm6mXyatmwib2JEfeaX0YBbSn7/ywA6iEo3lqJAPbTZHHqwN7RtK11wLnuDPZa60FkOiGoZdcMxPY1LQppzJveZ9mD6/84ZqK81POHDIP98AxDLt4gOFT5SD58pE5sOnFSFeRTx6QiAC4EcGbtJWJgksdZtIYEEenBgKMS4yn5zEbLczdq+qzjjG8Nddkva6p/JOQAZk5TErE1m4/3YbQp06XvR0jSKlRVegszjxT2KHVvkhudD3/76Upsfiu8TYuy3/GtHKGFkSjqeu37/SpyerlzU9Z+qnsibE0tK2PVqnUGcYuiMhiB9DWjg44sIF68DS2KBIfE1MuaK9ekzIzEhJMRS/T/ngBrBSyFWqqXrXxErRC5/aEGoZWwsM6VsOloji+gyrpOPNK21kE+r++yZdrCex0FhttCVfuZv/B3rvSA7y00noQ/GlT6kI9taS2TEYfWryrJDvTD48CJMWb9ZIATeurTl9oWG2lu4KStRJJeosq88no7cKBFfSLyrSflvzGc7I3d3ibLz6s0dXVw6yvWi0GTopQqliG1hRvcW3fjbBHeK+wrSPTnQgfw==.1; A=1l742shfnu3kk&t=1602887902&u=1602887902&sd=B&v=1&d=C08MKWPROjfMWNxR/l1QYQILtAZBKeiWCesjTWU75JTGK7eMXMJ9rXrr8PZ7qykr9cVizI87ukEvrmtNFtBKMI2a+s+jXgtuaWgyRaetfaPyhPCaZDS9722qS94qwtnB8pm6G3IS8Q5Yg+1ac4rZcasnkCjsPDJFcrTSkneHfZqP4ruq/0; XA=1l742shfnu3kk&t=1602887902&u=1602887902&sd=B&v=1&d=C08MKWPROjfMWNxR/l1QYQILtAZBKeiWCesjTWU75JTGK7eMXMJ9rXrr8PZ7qykr9cVizI87ukEvrmtNFtBKMI2a+s+jXgtuaWgyRaetfaPyhPCaZDS9722qS94qwtnB8pm6G3IS8Q5Yg+1ac4rZcasnkCjsPDJFcrTSkneHfZqP4ruq/0; JV=oEQkTDNL.wWMMdTExiIPrmzIcZS0FoECE4ExWL_lbnRV9SaM2iMj4AKRWfPujn3a9j7ZxhduvnzYHAZy4cTGCKRNQfPcRERaaZ0AcFaTOrIHCKB7U00vi73vSn2BJE06CFAC_lIq74u0PFJzCVlwCav97qUfUXQIjj4L6foeYFCGLvLkLfPREFnGRiMiZgw8i_vZD8L.CWQPO64I78fu1bns1ZuSklXssGHzfL730LsN2usvx5T88fcyzmYraWMT.Due1NnKe4NkyWj0XK8cg2.87wWT1ABxB4bo0JkEXNYWmIwdSNYQu4kgawsDW0kv3.ppBCUR3NNNsu68tBds2twg5Bftj0swR05OHkTtnZZ_hOtikltjUVkutl6.4vZBSRJGQJ7m8YxMcD6Bmj_uRuctB9tXEBI4bxXL1HE3g7f9dT8k.A54nMOiPDKpljwhKCB16yPoRpldtVyc9Ki6IIWkOb3s13ro1Inir6Mnt0L8NMDFSNefnjBTQ_MjKkUT1nKla8BaFKHMf3WhTsAflPZUJZrtKlk9kNncKuSZi.m903tgoba45OkoKtrRzP.b2v8Q3XaxkPGhAqmtd5dUxXTGGyXkqAkxINaHYv98O4WcHl.ILCbRWNxCdxC4wLcf2jqQ_SEbEMT_7voQVE_TXkWUQ2cbU2Y00arsrdSG7faF.EQDWpmFmRrOfV7sitSQdB290W6V0x9bUcz85kwh5Uv7nnU3wyfr8EN__l2sJsn2b05_nDUXX.81VUn0PsbT9Ak7rBZwZf9AND1ARwgastJGF0ZSpmtZozpdmHfMAEcnBuqTW9aSvellPqd2xxqwZInfZyIpvRcqvCPkpEMSSm_Te29cAGiAXdqBL_JdIqNU5pYDgakpBH2BtZPrduBWEdIw9PNwH2hpUxzcwHlGUZFlM7zWzjxrGYsQt.Mql8faImyf6ZP8baSagsaFR3HEl0cuPS7VWK_2Le02Q3_azv5Fo4rN829B8hKSl39Qb_ntBWPag93UUPVJXIMY; btpdb.2wzBV9u.dGZjLjEwNzQwOTQ2MA=REFZUw; btpdb.2wzBV9u.dGZjLjEwNzU0MTkzNg=REFZUw; btpdb.2wzBV9u.dGZjLjE0NDcxNDU=UkVRVUVTVFMuMjE; __gads=ID=7444c0ec18d52206:T=1604977013:S=ALNI_MYpQPlB5FkVm3KFHxq7QfkyZBdqXQ",
			},
			map[string][]string{
				"authority": {"www.yahoo.co.jp"},
				"cache-control": {"max-age=0"},
				"upgrade-insecure-requests": {"1"},
				"user-agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36"},
				"accept": {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"sec-fetch-site": {"none"},
				"sec-fetch-mode": {"navigate"},
				"sec-fetch-user": {"?1"},
				"sec-fetch-dest": {"document"},
				"accept-language": {"ja,en-US;q=0.9,en;q=0.8"},
				"cookie": {"F=a=AjGEC_IMvScXQ6jRpdCiF1Zhsv9NxTEB3jReaRXcj4fC9mnYtTgO5u5ZLB8FAmstQIK8ub4-&b=pZ98; B=1l742shfnu3kk&b=4&d=f.PR2GlpYF0ffQzA9.sN9PkjdO1k6F_IfPuvDYwl&s=0k&i=16bahJ2QKNr6oOzXHAKD; XB=1l742shfnu3kk&b=4&d=f.PR2GlpYF0ffQzA9.sN9PkjdO1k6F_IfPuvDYwl&s=0k&i=16bahJ2QKNr6oOzXHAKD; Y=v=1&n=bj8bbcaff7t61&l=c0c0e_o/o&p=m2kvvjp052000000&ig=01767&r=115&lg=ja-JP&intl=jp; T=z=BA4gfBBoGqfBIXJK3.cWRIkMDc3MwYxNjFOMDUwTzA-&sk=DAAwvb7OXbcH80&ks=EAAHuBJ4plFljvAGaxRTAEuEg--~F&kt=EAAD0Pb3pbZ.h4xxZ8GbOgmbw--~E&ku=FAAMEYCIQCrcvvQ2x.iLqbg_A1F25vdOsdlKn7vNCl4CPRA.GjH7QIhAIWRpIOGK2H8R1fW7g_LLf8mE9ekXSPP8FIvBByI7_Y3~B&d=dGlwATJGeDBEQwFhAVlBRQFnATRTUVJLS1U0WERWRVhWU0ZENFBVNFNCWE5BAXNsAU56QXdOQUUyTVRZNU56STNPRGMtAXNjAXd3dwF6egFCQTRnZkJBMko-; SSL=v=1&s=H8a9563jCll2MrcDpijrxS8QQ2tVAOJnetfb4_UMEiEdQThv5tHTfRUHBhp18zgi.wAgwzbYeLgKffvW9yJXlw--&kv=0; _n=DPFnxL-4PtCUH4M4G1idsx8MuEQMPWcLNUO8kFSRYa9MTwBAZjRLGMs4GnGMVjk4IEKcOFHNMdXMF-goe3DvrEbYbvnErNZIzgZtBdAP-FIPPE1O8QeolMcT6Al37amXltSx7tR7011AcqLTsPE7yDPkgcrGKPwBubo6vjDJp2GYHPfkDaWhTRyWowvzKG7b3GAs38apZuiUCk4chVfNSJZoB_IeuJ4rohz3m6aYWSm2NWiKe-0xwaDJebgmhI7e16ab3lu8TLqR2wzqTq0j_jIXMEl49l0bL8fBeH87gYvXcsjvape1d3l3QeVRb1SwgfVCLj3QJKAvowjUydSGYBUemL2j1Rrb6gm_Cbf4Vb8kWDJVcE1abCPQG-N5uPobJA1DEl623ihTOfSePvJIAibHDs2MUfAklkAVAqO530Ha66SBl1lZZRJIxy1IFigSWTLknJHEuvj981Qehwv8E63r6KyinP89GpN_fyMh1h8YVki6NQhvp0cbN8jfc9yPHleV-9t8IpJHDZzUKS3Db6NNDDZdo1LhM63xrCxbDMpdeNFtahRkA2KfwmSMuZMu0CO2x-pwHAhgLoUUFP_Yik4hmY5gGnQZasrsTk_2ZJ8U4fLhyGcgnxhfZtvPMz5qIYKwHDMOm6vf5C_ql9-HJjR_zedsb5IkIVf2zE1ZVBKUhjLwrmJNof93Ag2dH3h4CFRTaGlfoxNlWJe0S5yuNnp6YUxqMxAIRf6HwXCnKdKGD07zblPPqvte9wALQpOmf6BwmpaFpqb-VP1-qfj3WVv2gX7HAa_G75mFVush-mtP45kSO4TCBG9j9MutYOHCkjYFw3V3HCY9UaLeUY56mHrK63ySJPuMXrX7_lKJemYcbvId4TsEo_kuufGWM-RueWM8EHvpvS8UooQveOwRN15hEJL-Zdi_8JHZOlCegHovMbJWLY3avtVJbmip_oV_p6LNlSWkeM8E5RyBQuPwZskXjSX21Sv2gOXSCVmyNI-RYWlIEnd2tIfMaU5gGyUnMBJTIqDIP-nsr9AjkEslJAJPUCX_mMtRM8_EGhoX1OS8m_jqDNcJBzR1dE1KFN8du8Xb1cEfPHRy8V2mHHC8FavMlgS9esUgzELkU6vQYnVs_-YGJsP1osFFx5ebk_c8.2; XC=JnDKyz8lyzvMSPDcpm99NAG3Cm6mXyatmwib2JEfeaX0YBbSn7/ywA6iEo3lqJAPbTZHHqwN7RtK11wLnuDPZa60FkOiGoZdcMxPY1LQppzJveZ9mD6/84ZqK81POHDIP98AxDLt4gOFT5SD58pE5sOnFSFeRTx6QiAC4EcGbtJWJgksdZtIYEEenBgKMS4yn5zEbLczdq+qzjjG8Nddkva6p/JOQAZk5TErE1m4/3YbQp06XvR0jSKlRVegszjxT2KHVvkhudD3/76Upsfiu8TYuy3/GtHKGFkSjqeu37/SpyerlzU9Z+qnsibE0tK2PVqnUGcYuiMhiB9DWjg44sIF68DS2KBIfE1MuaK9ekzIzEhJMRS/T/ngBrBSyFWqqXrXxErRC5/aEGoZWwsM6VsOloji+gyrpOPNK21kE+r++yZdrCex0FhttCVfuZv/B3rvSA7y00noQ/GlT6kI9taS2TEYfWryrJDvTD48CJMWb9ZIATeurTl9oWG2lu4KStRJJeosq88no7cKBFfSLyrSflvzGc7I3d3ibLz6s0dXVw6yvWi0GTopQqliG1hRvcW3fjbBHeK+wrSPTnQgfw==.1; A=1l742shfnu3kk&t=1602887902&u=1602887902&sd=B&v=1&d=C08MKWPROjfMWNxR/l1QYQILtAZBKeiWCesjTWU75JTGK7eMXMJ9rXrr8PZ7qykr9cVizI87ukEvrmtNFtBKMI2a+s+jXgtuaWgyRaetfaPyhPCaZDS9722qS94qwtnB8pm6G3IS8Q5Yg+1ac4rZcasnkCjsPDJFcrTSkneHfZqP4ruq/0; XA=1l742shfnu3kk&t=1602887902&u=1602887902&sd=B&v=1&d=C08MKWPROjfMWNxR/l1QYQILtAZBKeiWCesjTWU75JTGK7eMXMJ9rXrr8PZ7qykr9cVizI87ukEvrmtNFtBKMI2a+s+jXgtuaWgyRaetfaPyhPCaZDS9722qS94qwtnB8pm6G3IS8Q5Yg+1ac4rZcasnkCjsPDJFcrTSkneHfZqP4ruq/0; JV=oEQkTDNL.wWMMdTExiIPrmzIcZS0FoECE4ExWL_lbnRV9SaM2iMj4AKRWfPujn3a9j7ZxhduvnzYHAZy4cTGCKRNQfPcRERaaZ0AcFaTOrIHCKB7U00vi73vSn2BJE06CFAC_lIq74u0PFJzCVlwCav97qUfUXQIjj4L6foeYFCGLvLkLfPREFnGRiMiZgw8i_vZD8L.CWQPO64I78fu1bns1ZuSklXssGHzfL730LsN2usvx5T88fcyzmYraWMT.Due1NnKe4NkyWj0XK8cg2.87wWT1ABxB4bo0JkEXNYWmIwdSNYQu4kgawsDW0kv3.ppBCUR3NNNsu68tBds2twg5Bftj0swR05OHkTtnZZ_hOtikltjUVkutl6.4vZBSRJGQJ7m8YxMcD6Bmj_uRuctB9tXEBI4bxXL1HE3g7f9dT8k.A54nMOiPDKpljwhKCB16yPoRpldtVyc9Ki6IIWkOb3s13ro1Inir6Mnt0L8NMDFSNefnjBTQ_MjKkUT1nKla8BaFKHMf3WhTsAflPZUJZrtKlk9kNncKuSZi.m903tgoba45OkoKtrRzP.b2v8Q3XaxkPGhAqmtd5dUxXTGGyXkqAkxINaHYv98O4WcHl.ILCbRWNxCdxC4wLcf2jqQ_SEbEMT_7voQVE_TXkWUQ2cbU2Y00arsrdSG7faF.EQDWpmFmRrOfV7sitSQdB290W6V0x9bUcz85kwh5Uv7nnU3wyfr8EN__l2sJsn2b05_nDUXX.81VUn0PsbT9Ak7rBZwZf9AND1ARwgastJGF0ZSpmtZozpdmHfMAEcnBuqTW9aSvellPqd2xxqwZInfZyIpvRcqvCPkpEMSSm_Te29cAGiAXdqBL_JdIqNU5pYDgakpBH2BtZPrduBWEdIw9PNwH2hpUxzcwHlGUZFlM7zWzjxrGYsQt.Mql8faImyf6ZP8baSagsaFR3HEl0cuPS7VWK_2Le02Q3_azv5Fo4rN829B8hKSl39Qb_ntBWPag93UUPVJXIMY; btpdb.2wzBV9u.dGZjLjEwNzQwOTQ2MA=REFZUw; btpdb.2wzBV9u.dGZjLjEwNzU0MTkzNg=REFZUw; btpdb.2wzBV9u.dGZjLjE0NDcxNDU=UkVRVUVTVFMuMjE; __gads=ID=7444c0ec18d52206:T=1604977013:S=ALNI_MYpQPlB5FkVm3KFHxq7QfkyZBdqXQ"},
			},
		},
		{
			[]string{"abc"},
			map[string][]string{},
		},
		{
			[]string{},
			map[string][]string{},
		},
	} {
		a := NewAttacker(&http.Client{}, &url, &method, test.headers, &duration, &rate, &workerNumber)
		h := a.parseHeader()
		if len(h) != len(test.expected) {
			t.Errorf("[TestParseHeader] header length : actual = %d, expected = %d", len(h), len(test.expected))
		}
	}
}

func TestShouldSendSignal(t *testing.T) {
	url := ""
	method := ""
	var headers []string
	workerNumber := 0
	for _, test := range []struct {
		name string
		rate int
		duration int
		rc int
		elapsed time.Duration
		expected bool
	} {
		{"RateIsNotMet", 10, 10, 99, 1 * time.Second, true},
		{"RateIsMet", 10, 10, 100, 1 * time.Second, false},
		{"DurationHasPassed", 0, 10, 1, 1 * time.Second, true},
		{"DurationHasNotPassed", 0, 10, 1, 11 * time.Second, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			a := NewAttacker(&http.Client{}, &url, &method, headers, &test.duration, &test.rate, &workerNumber)
			res := a.shouldSendSignal(&test.rc, &test.elapsed)
			if res != test.expected {
				t.Errorf("[TestShouldSendSignal] failed : actual = %t, expected = %t", res, test.expected)
			}
		})
	}
}