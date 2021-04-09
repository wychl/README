package main

import (
	networking "istio.io/api/networking/v1alpha3"
)

func convertTOStringMatch(way, value string) *networking.StringMatch {
	sm := &networking.StringMatch{}
	switch way {
	case prefix:
		sm.MatchType = &networking.StringMatch_Prefix{
			Prefix: value,
		}
	case exact:
		sm.MatchType = &networking.StringMatch_Exact{
			Exact: value,
		}
	case regex:
		sm.MatchType = &networking.StringMatch_Regex{
			Regex: value,
		}

	}

	return sm
}
