package ai

import (
	"fmt"
	"github.com/jasonlovesdoggo/roastme/internal/analysis"
	"math/rand"
)

// generateLocalRoast generates a roast without using an external AI service
func generateLocalRoast(patterns analysis.CommandPattern, complexity ComplexityLevel) string {
	var basicRoasts []string

	if len(patterns.RepeatedCommands) > 0 {
		cmd := patterns.RepeatedCommands[0].Command
		count := patterns.RepeatedCommands[0].Count
		basicRoasts = append(basicRoasts, fmt.Sprintf("I see you've used '%s' %d times. Having memory issues or just really, really in love with that command?", cmd, count))
	}

	if len(patterns.FailedCommands) > 0 {
		basicRoasts = append(basicRoasts, fmt.Sprintf("Nice typos! Maybe typing lessons should be in your future before attempting '%s'.", patterns.FailedCommands[0]))
	}

	if len(patterns.ComplexCommands) > 0 {
		basicRoasts = append(basicRoasts, "Wow, those complex commands! Trying to impress an invisible audience or just afraid of using separate lines?")
	}

	if patterns.Indecisive {
		basicRoasts = append(basicRoasts, "All those cd's and ls's... are you exploring your filesystem or just completely lost in there?")
	}

	if len(patterns.TimeWasters) > 0 {
		waster := patterns.TimeWasters[0]
		basicRoasts = append(basicRoasts, fmt.Sprintf("I see you're visiting %s. Working hard or hardly working, eh?", waster))
	}

	// Add skill level roasts
	skill := patterns.SkillLevel
	if skill == "beginner" {
		basicRoasts = append(basicRoasts, "Your command history screams 'I just discovered what a terminal is'. How adorable.")
	} else if skill == "intermediate" {
		basicRoasts = append(basicRoasts, "Your command history is like a mediocre pizza - it gets the job done but nobody's impressed.")
	} else {
		basicRoasts = append(basicRoasts, "Fancy commands! Overcompensating for something or just showing off to nobody?")
	}

	// Make sure we have at least one roast
	if len(basicRoasts) == 0 {
		basicRoasts = append(basicRoasts, "I can't even roast your command history - it's that boring. Try doing something interesting first!")
	}

	switch complexity {
	case SimpleRoast:
		// Just return a single roast
		return basicRoasts[rand.Intn(len(basicRoasts))]

	case NormalRoast:
		// Combine 2 roasts
		if len(basicRoasts) >= 2 {
			idx1 := rand.Intn(len(basicRoasts))
			idx2 := (idx1 + 1 + rand.Intn(len(basicRoasts)-1)) % len(basicRoasts)
			return basicRoasts[idx1] + " " + basicRoasts[idx2]
		}
		return basicRoasts[0]

	case ComplexRoast:
		// Build a more complex roast with intro and conclusion
		result := "Well, well, well. Looking at your command history is like reading a diary of technical confusion.\n\n"

		// Add 2-3 specific roasts
		numRoasts := min(len(basicRoasts), 3)
		usedIndices := map[int]bool{}

		for i := 0; i < numRoasts; i++ {
			idx := rand.Intn(len(basicRoasts))
			for usedIndices[idx] {
				idx = rand.Intn(len(basicRoasts))
			}
			usedIndices[idx] = true
			result += basicRoasts[idx] + " "
		}

		result += "\n\nMaybe someday you'll graduate to actually knowing what you're doing in the terminal, but today is clearly not that day."
		return result

	case BrutalRoast:
		// Build a comprehensive multi-paragraph roast
		result := "I've seen some sad terminal histories in my time, but yours takes the award for 'Most Likely To Make Linus Torvalds Weep.'\n\n"

		// Use all available roasts
		for i, roast := range basicRoasts {
			result += roast + " "
			if (i+1)%2 == 0 {
				result += "\n\n"
			}
		}

		result += "\n\nIf your code is anything like your command history, I'm guessing Stack Overflow isn't just a website for you - it's a lifeline. Your terminal doesn't need a history feature; it needs therapy after what you've put it through."

		result += "\n\nThe good news? Even trained monkeys eventually learn which buttons not to push. There's hope for you yet."

		return result
	}

	// Default fallback
	return basicRoasts[rand.Intn(len(basicRoasts))]
}
