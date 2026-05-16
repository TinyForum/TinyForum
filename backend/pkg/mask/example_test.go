package mask_test

import (
	"fmt"
	"tiny-forum/pkg/mask"
)

type Address struct {
	Detail string `mask:"address,keep=4"`
	City   string `mask:"full"`
}

type User struct {
	Name     string   `mask:"name"`
	Phone    string   `mask:"mobile"`
	Email    string   `mask:"email"`
	IDCard   string   `mask:"idcard"`
	BankCard string   `mask:"bankcard"`
	Password string   `mask:"full"`
	Addr     Address  `mask:"-"`
	Tags     []string `mask:"full"`
	Nick     *string  `mask:"name"`
	Ignore   string   `mask:"-"`
}

func ExampleMask() {
	nick := "小欧阳"
	u := &User{
		Name:     "欧阳锋",
		Phone:    "13812345678",
		Email:    "feng.ouyang@example.com",
		IDCard:   "11010119900307663X",
		BankCard: "6222021234567890123",
		Password: "mySecret123",
		Addr: Address{
			Detail: "上海市黄浦区南京东路100号",
			City:   "上海",
		},
		Tags:   []string{"gopher", "runner"},
		Nick:   &nick,
		Ignore: "should not show",
	}

	if err := mask.Mask(u); err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s, Phone: %s, Email: %s, IDCard: %s, BankCard: %s, Password: %s, Addr: %+v, Tags: %v, Nick: %s, Ignore: %s\n",
		u.Name, u.Phone, u.Email, u.IDCard, u.BankCard, u.Password, u.Addr, u.Tags, *u.Nick, u.Ignore)
	// Output:
	// Name: 欧**, Phone: 138****5678, Email: f***@example.com, IDCard: 110101**********X, BankCard: 622202***********0123, Password: ********, Addr: {Detail:上海**** City:**}, Tags: [**** ****], Nick: 小**, Ignore: should not show
}

func ExampleMaskCopy() {
	original := &User{Name: "张三", Phone: "13912345678"}
	copied, _ := mask.MaskCopy(original)
	// original 不变
	fmt.Printf("original: %+v\n", original)
	fmt.Printf("masked: %+v\n", copied)
}
