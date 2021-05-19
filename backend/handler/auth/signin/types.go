package signin

type SigninInput struct {
	Email    string
	Password string
}

type SigninOutput struct {
	AccessToken string `json:"accessToken"`
}

type Mutation struct {
	Signin *SigninOutput
}

type SigninArgs struct {
	Input SigninInput
}
