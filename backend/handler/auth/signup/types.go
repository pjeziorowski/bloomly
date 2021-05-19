package signup

type SignupInput struct {
	Email    string
	Password string
}

type SignupOutput struct {
	AccessToken string `json:"accessToken"`
}

type Mutation struct {
	Signup *SignupOutput
}

type SignupArgs struct {
	Arg1 SignupInput
}
