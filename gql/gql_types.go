package gql

import (
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"
)

func GenCompInfoToDBCompInfo(gen *generated.CompatibilityInfoInput) *util.CompatibilityInfo {
	if gen == nil {
		return nil
	}
	return &util.CompatibilityInfo{
		Ea:         GenCompToDBComp(gen.Ea),
		Exp:        GenCompToDBComp(gen.Exp),
		Controller: GenControllerCompToDBControllerComp(gen.Controller),
	}
}

func GenCompToDBComp(gen *generated.CompatibilityInput) util.Compatibility {
	r := util.Compatibility{
		State: string(gen.State),
	}
	SetINN(gen.Note, &r.Note)
	return r
}

func GenControllerCompToDBControllerComp(gen *generated.ControllerCompatibilityInput) util.Compatibility {
	r := util.Compatibility{
		State: string(gen.State),
	}
	SetINN(gen.Note, &r.Note)
	return r
}

func GenAIDisclosureInfoToDBAIDisclosureInfo(gen *generated.AIUseDisclosureInput) *util.AIUseDisclosureInfo {
	if gen == nil {
		return nil
	}

	//DisclosureType AIUseDisclosureType
	//DisclosureString string
	/*
	return &util.AIUseDisclosureInfo{
		DisclosureType:         GenAIDTToDBAIDT(gen.DisclosureType),
		DisclosureString:        gen.DisclosureString,
	}
*/
	r := &util.AIUseDisclosureInfo{
		DisclosureType: string(gen.DisclosureType),
	}
	SetINN(gen.DisclosureString, &r.DisclosureString)
	return r
}