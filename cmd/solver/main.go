package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/osmanyasin/grand-prix/internal/io"
	"github.com/osmanyasin/grand-prix/internal/models"
	"github.com/osmanyasin/grand-prix/internal/solver"
)

func main() {
	levelPath := flag.String("level", "levels/L4.txt", "Path to the level JSON file")
	outPath := flag.String("out", "strategies/greedy_strategy.txt", "Path to save the generated strategy")
	flag.Parse()

	fmt.Printf("Loading level: %s\n", *levelPath)

	// 1. Read the level configuration
	config, err := io.ReadLevel[models.LevelConfig](*levelPath)
	if err != nil {
		log.Fatalf("Failed to load level config: %v", err)
	}

	// 2. Run the Greedy Algorithm
	fmt.Println("Running Greedy Optimizer...")
	strategy := solver.GenerateGreedyStrategy(config)

	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(*outPath), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// 3. Save the result
	err = io.WriteStrategy(*outPath, strategy)

	// 3. Save the result
	err = io.WriteStrategy(*outPath, strategy) // Ensure this exists in your io package
	if err != nil {
		log.Fatalf("Failed to save strategy: %v", err)
	}

	fmt.Printf("Successfully generated strategy at: %s\n", *outPath)
	fmt.Printf("Run the evaluator to test it:\n")

	// Print a helpful hint for the user
	levelNum := filepath.Base(*levelPath)[1:2]
	fmt.Printf("go run ./cmd/evaluator -level %s -strategy %s -levelnum %s\n", *levelPath, *outPath, levelNum)
}
