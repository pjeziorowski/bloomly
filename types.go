package main

type SignupInput struct {
	Email    string
	Password string
}

type SigninInput struct {
	Email    string
	Password string
}

type SignupOutput struct {
	AccessToken string
}

type SigninOutput struct {
	AccessToken string
}
