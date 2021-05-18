package main

type SigninArgs struct {
	Arg1 SigninInput
}

type SignupArgs struct {
	Arg1 SignupInput
}

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
