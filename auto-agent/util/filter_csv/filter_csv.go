package main

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/toisin/astro-world/auto-agent/util"
)

const promptIDIdx = 2

var promptIDs = []string{
	"two_records.same_record.2.1",
	"single_record.hide_performance.1",
	"single_record.show_performance.1",
	"single_record.think_enough_info.2nd_record.show_performance.1.1.1",
	"two_records.hide_performance.2.2",
	"two_records.show_performance.2.2",
	"causal.target_nonvarying.Q_why.2.2.1.1",
	"causal.uncontrolled.Q_why.2.2.1.2",
	"causal.controlled.Q_why.2.2.1.3",
	"noncausal.target_nonvarying.Q_why.2.2.2.1",
	"noncausal.uncontrolled.Q_why.2.2.2.2",
	"causal.Q_why.1",
	"non-causal.Q_why.2",
	"help.Q_why.3.1",
	"wrong.wrong.1.2",
	"wrong.correct.1.1",
	"causal.controlled.correct.sure.Q_someone_disagree.2.2.1.3.2.1                                        causal.controlled.correct.challenge.2.2.1.3.2",
	"noncausal.controlled.correct.Q_someone_disagree.2.2.2.3.2",
	"non-causal.correct.correct.Q_someone_disagree.2.3.2",
	"causal.correct.correct.Q_someone_disagree.1.1",
	"two_records.MC_is_causal.2.2",
	"unsure.target_varying.2.2.3.2",
	"unsure.dont_know.2.2.3.2.3",
	"causal.controlled.correct.unsure.2.2.1.3.2.2",
	"causal.uncontrolled.MC_are_sure.2.2.1.2",
	"noncausal.controlled.correct.MC_are_sure.2.2.2.3.2",
	"causal.controlled.correct.MC_are_sure.2.2.1.3.2",
	"single_record.MC_enough_info.1",
	"single_record.think_enough_info.Q_are_sure.1.1",
	"causal.uncontrolled.MC_are_sure.2.2.1.2",
	"causal.controlled.correct.MC_are_sure.2.2.1.3.2",
	"noncausal.uncontrolled.MC_are_sure.2.2.2.2",
	"noncausal.controlled.correct.MC_are_sure.2.2.2.3.2",
	"start_prediction",
	"wrong_factor.correct.correct_factors.prediction.1.2.2",
	"pick_team",
}

func main() {
	f := util.OpenFileFromArg()
	defer f.Close()

	r := util.NewCSVReader(f)
	w := csv.NewWriter(os.Stdout)

	// Header
	header, err := r.Read()
	util.MaybeExit(err)

	err = w.Write(header)
	util.MaybeExit(err)

	set := map[string]struct{}{}
	for _, word := range promptIDs {
		set[word] = struct{}{}
	}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		util.MaybeExit(err)

		promptID := row[promptIDIdx]
		if _, ok := set[promptID]; ok {
			err = w.Write(row)
			util.MaybeExit(err)
		}
	}

	w.Flush()
}
