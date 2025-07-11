package ontology

// ContextType represents the type of context for tag classification
type ContextType string

const (
	ContextTypeColleges             ContextType = "Colleges"
	ContextTypeCertifications       ContextType = "Certifications"
	ContextTypeEntranceExams        ContextType = "EntranceExams"
	ContextTypeAPExams              ContextType = "APExams"
	ContextTypeUserGeneratedContent ContextType = "UserGeneratedContent"
	ContextTypeDoD                  ContextType = "DoD"
	ContextTypeNone                 ContextType = "None"
)

// TagType represents the type of a tag in the tree structure
type TagType string

const (
	TagTypeCategory         TagType = "Category"
	TagTypeSubCategory      TagType = "SubCategory"
	TagTypeUniversity       TagType = "University"
	TagTypeRegion           TagType = "Region"
	TagTypeDepartment       TagType = "Department"
	TagTypeCourse           TagType = "Course"
	TagTypeTopic            TagType = "Topic"
	TagTypeUserFolder       TagType = "UserFolder"
	TagTypeUserTopic        TagType = "UserTopic"
	TagTypeCertifyingAgency TagType = "Certifying_Agency"
	TagTypeCertification    TagType = "Certification"
	TagTypeDomain           TagType = "Domain"
	TagTypeModule           TagType = "Module"
	TagTypeEntranceExam     TagType = "Entrance_Exam"
	TagTypeAPExam           TagType = "AP_Exam"
	TagTypeUserContent      TagType = "UserContent"
	TagTypeBranch           TagType = "Branch"
	TagTypeInstructionType  TagType = "Instruction_Type"
	TagTypeInstructionGroup TagType = "Instruction_Group"
	TagTypeInstruction      TagType = "Instruction"
	TagTypeChapter          TagType = "Chapter"
	TagTypeSection          TagType = "Section"
	TagTypePart             TagType = "Part"
	TagTypeNone             TagType = "None"
)

// TagOntology defines the mapping between context type, depth, and tag types
type TagOntology struct {
	ContextType  ContextType
	HeaderLength int
	TagTypes     []TagType
}
