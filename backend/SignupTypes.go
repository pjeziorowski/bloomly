package main

type SampleInput struct {
	Username string
	Password string
}

type SignupInput struct {
	Email    string
	Password string
}

type SampleOutput struct {
	AccessToken string
}

type SignupOutput struct {
	Email       string
	AccessToken string
}

type Mutation struct {
	Signup *SampleOutput
}

type SignupArgs struct {
	Arg1 SampleInput
}
