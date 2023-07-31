package di

import "go.uber.org/dig"

var Container = dig.New()

func MustProvide(constructor interface{}, opts ...dig.ProvideOption) {
	err := Container.Provide(constructor, opts...)
	if err != nil {
		panic(err)
	}
}

func MustInvoke(function interface{}, opts ...dig.InvokeOption) {
	err := Container.Invoke(function, opts...)
	if err != nil {
		panic(err)
	}
}
