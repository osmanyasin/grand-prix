package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/osmanyasin/grand-prix/internal/io"
	"github.com/osmanyasin/grand-prix/internal/models"
	"github.com/osmanyasin/grand-prix/internal/scoring"
	"github.com/osmanyasin/grand-prix/internal/simulator"
)

func main() {
	levelPath := flag.String("level", "", "Path to the level JSON file")
	strategyPath := flag.String("strategy", "", "Path to the strategy JSON/txt file")
	levelNum := flag.Int("levelnum", 4, "The level number being evaluated (1-4)")
	flag.Parse()

	if *levelPath == "" || *strategyPath == "" {
		log.Fatal("Both -level and -strategy flags are required.")
	}

	config, err := io.ReadLevel[models.LevelConfig](*levelPath)
	if err != nil {
		log.Fatalf("Failed to load level config: %v", err)
	}

	strategy, err := io.ReadLevel[models.Strategy](*strategyPath)
	if err != nil {
		log.Fatalf("Failed to load strategy: %v", err)
	}

	fmt.Println("Starting race simulation...")
	finalState, err := simulator.EvaluateRace(config, strategy)
	if err != nil {
		log.Fatalf("Simulation failed: %v", err)
	}

	final, base, fuel, tyre := scoring.CalculateFinalScore(*levelNum, finalState, config)

	fmt.Println("\n--- Race Results ---")
	fmt.Printf("Total Time:          %.3f s\n", finalState.TotalTimeSeconds)
	fmt.Printf("Total Fuel Used:     %.3f L\n", finalState.TotalFuelUsedLitres)
	fmt.Printf("Total Tyre Wear:     %.3f\n", finalState.TotalTyreDegradation)
	fmt.Printf("Blowouts:            %d\n", finalState.NumberOfBlowouts)

	if finalState.IsLimpMode {
		fmt.Println("STATUS:              FINISHED IN LIMP MODE")
	}

	fmt.Println("\n--- Final Scoring ---")
	fmt.Printf("Base Score:          + %.0f\n", base)
	if *levelNum >= 2 {
		fmt.Printf("Fuel Bonus:          + %.0f\n", fuel)
	}
	if *levelNum >= 4 {
		fmt.Printf("Tyre Bonus:          + %.0f\n", tyre)
	}
	fmt.Println("-----------------------")
	fmt.Printf("FINAL SCORE:         %.0f\n\n", final)
}
