package main

import (
	"fmt"
	"log"

	ris "ris-go/lib"
)

func main() {
	// load data
	maleData, err := ris.LoadDataFromCSV("data/male.csv")
	if err != nil {
		log.Fatal(err)
	}
	femaleData, err := ris.LoadDataFromCSV("data/female.csv")
	if err != nil {
		log.Fatal(err)
	}

	// calculate fitting parameters
	maleFit, err := ris.FitRISParamsScipy(maleData, 100)
	if err != nil {
		log.Fatal(err)
	}

	// calculate fitting parameters
	femaleFit, err := ris.FitRISParamsNelder(femaleData, 100)
	if err != nil {
		log.Fatal(err)
	}

	//  parameters for classic
	fmt.Println("=== 2-LIFT STREETLIFTING PARAMETERS ===")
	fmt.Println("\nMALE PARAMETERS:")
	fmt.Printf("A: %.5f\n", maleFit.Params.A)
	fmt.Printf("K: %.5f\n", maleFit.Params.K)
	fmt.Printf("B: %.5f\n", maleFit.Params.B)
	fmt.Printf("v: %.5f\n", maleFit.Params.V)
	fmt.Printf("Q: %.5f\n", maleFit.Params.Q)
	fmt.Printf("RMSE: %.2f kg\n", maleFit.RMSE)

	fmt.Println("\nFEMALE PARAMETERS:")
	fmt.Printf("A: %.5f\n", femaleFit.Params.A)
	fmt.Printf("K: %.5f\n", femaleFit.Params.K)
	fmt.Printf("B: %.5f\n", femaleFit.Params.B)
	fmt.Printf("v: %.5f\n", femaleFit.Params.V)
	fmt.Printf("Q: %.5f\n", femaleFit.Params.Q)
	fmt.Printf("RMSE: %.2f kg\n", femaleFit.RMSE)
}
