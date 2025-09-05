package ontology

var tagOntology = []TagOntology{
	// AP Exams
	{ContextTypeAPExams, 3, []TagType{
		TagTypeCategory,
		TagTypeAPExam,
		TagTypeTopic,
	}},
	{ContextTypeAPExams, 4, []TagType{
		TagTypeCategory,
		TagTypeAPExam,
		TagTypeModule,
		TagTypeTopic,
	}},
	{ContextTypeAPExams, 5, []TagType{
		TagTypeCategory,
		TagTypeAPExam,
		TagTypeModule,
		TagTypeModule,
		TagTypeTopic,
	}},
	{ContextTypeAPExams, 6, []TagType{
		TagTypeCategory,
		TagTypeAPExam,
		TagTypeModule,
		TagTypeModule,
		TagTypeModule,
		TagTypeTopic,
	}},

	// Certifications
	{ContextTypeCertifications, 4, []TagType{
		TagTypeCategory,
		TagTypeCertifyingAgency,
		TagTypeCertification,
		TagTypeTopic,
	}},
	{ContextTypeCertifications, 5, []TagType{
		TagTypeCategory,
		TagTypeCertifyingAgency,
		TagTypeCertification,
		TagTypeModule,
		TagTypeTopic,
	}},
	{ContextTypeCertifications, 6, []TagType{
		TagTypeCategory,
		TagTypeCertifyingAgency,
		TagTypeCertification,
		TagTypeDomain,
		TagTypeModule,
		TagTypeTopic,
	}},
	{ContextTypeCertifications, 7, []TagType{
		TagTypeCategory,
		TagTypeCertifyingAgency,
		TagTypeCertification,
		TagTypeDomain,
		TagTypeDomain,
		TagTypeModule,
		TagTypeTopic,
	}},
	{ContextTypeCertifications, 8, []TagType{
		TagTypeCategory,
		TagTypeCertifyingAgency,
		TagTypeCertification,
		TagTypeDomain,
		TagTypeDomain,
		TagTypeDomain,
		TagTypeModule,
		TagTypeTopic,
	}},
	{ContextTypeCertifications, 9, []TagType{
		TagTypeCategory,
		TagTypeCertifyingAgency,
		TagTypeCertification,
		TagTypeDomain,
		TagTypeDomain,
		TagTypeDomain,
		TagTypeDomain,
		TagTypeModule,
		TagTypeTopic,
	}},

	// Colleges
	{ContextTypeColleges, 6, []TagType{
		TagTypeCategory,
		TagTypeRegion,
		TagTypeUniversity,
		TagTypeDepartment,
		TagTypeCourse,
		TagTypeTopic,
	}},

	// DoD
	{ContextTypeDoD, 4, []TagType{
		TagTypeCategory,
		TagTypeBranch,
		TagTypeInstructionType,
		TagTypeInstructionGroup,
	}},
	{ContextTypeDoD, 5, []TagType{
		TagTypeCategory,
		TagTypeBranch,
		TagTypeInstructionType,
		TagTypeInstructionGroup,
		TagTypeInstruction,
	}},
	{ContextTypeDoD, 6, []TagType{
		TagTypeCategory,
		TagTypeBranch,
		TagTypeInstructionType,
		TagTypeInstructionGroup,
		TagTypeInstruction,
		TagTypeSection,
	}},
	{ContextTypeDoD, 7, []TagType{
		TagTypeCategory,
		TagTypeBranch,
		TagTypeInstructionType,
		TagTypeInstructionGroup,
		TagTypeInstruction,
		TagTypeSection,
		TagTypeChapter,
	}},
	{ContextTypeDoD, 8, []TagType{
		TagTypeCategory,
		TagTypeBranch,
		TagTypeInstructionType,
		TagTypeInstructionGroup,
		TagTypeInstruction,
		TagTypeSection,
		TagTypeChapter,
		TagTypePart,
	}},

	// Entrance Exams
	{ContextTypeEntranceExams, 3, []TagType{
		TagTypeCategory,
		TagTypeEntranceExam,
		TagTypeTopic,
	}},
	{ContextTypeEntranceExams, 4, []TagType{
		TagTypeCategory,
		TagTypeEntranceExam,
		TagTypeModule,
		TagTypeTopic,
	}},
	{ContextTypeEntranceExams, 5, []TagType{
		TagTypeCategory,
		TagTypeEntranceExam,
		TagTypeModule,
		TagTypeModule,
		TagTypeTopic,
	}},
	{ContextTypeEntranceExams, 6, []TagType{
		TagTypeCategory,
		TagTypeEntranceExam,
		TagTypeModule,
		TagTypeModule,
		TagTypeModule,
		TagTypeTopic,
	}},

	// Encyclopedia
	{ContextTypeEncyclopedia, 4, []TagType{
		TagTypeCategory,
		TagTypeVolume,
		TagTypeRange,
		TagTypeTopic,
	}},
}
