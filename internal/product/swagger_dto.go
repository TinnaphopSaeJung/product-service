package product

type PatchProductSwaggerRequest struct {
	Name        *string  `json:"name" example:"Gaming Keyboard"`
	Description *string  `json:"description" example:"Updated description"`
	SalePrice   *float64 `json:"sale_price" example:"1290"`
	Price       *float64 `json:"price" example:"1590"`
}

type CreateProductSuccessSwaggerResponse struct {
	Successful bool            `json:"successful" example:"true"`
	ErrorCode  *string         `json:"error_code" swaggertype:"string" example:"null"`
	Data       ProductResponse `json:"data"`
}

type CreateProductErrorSwaggerResponse struct {
	Successful bool    `json:"successful" example:"false"`
	ErrorCode  string  `json:"error_code" example:"VALIDATION_ERROR"`
	Data       *string `json:"data" swaggertype:"string" example:"null"`
}

type PatchProductSuccessSwaggerResponse struct {
	Successful bool    `json:"successful" example:"true"`
	ErrorCode  *string `json:"error_code" swaggertype:"string" example:"null"`
}

type PatchProductErrorSwaggerResponse struct {
	Successful bool   `json:"successful" example:"false"`
	ErrorCode  string `json:"error_code" example:"VALIDATION_ERROR"`
}
