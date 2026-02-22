package engine

import "github.com/careermanifest/backend/internal/dto"

// ============================================================
// CAREER-SPECIFIC ENRICHMENT DATA (India-focused, realistic)
// ============================================================

// GetSalaryProjection returns 5-year salary data for a career.
func GetSalaryProjection(career Career) dto.SalaryProjection {
	projections := map[Career]dto.SalaryProjection{
		CareerIT: {
			Year1: "₹4-8 LPA", Year2: "₹6-12 LPA", Year3: "₹10-18 LPA",
			Year4: "₹14-25 LPA", Year5: "₹18-35 LPA",
		},
		CareerMBA: {
			Year1: "₹8-15 LPA", Year2: "₹10-20 LPA", Year3: "₹14-28 LPA",
			Year4: "₹18-35 LPA", Year5: "₹22-50 LPA",
		},
		CareerGovt: {
			Year1: "₹5-8 LPA", Year2: "₹5.5-9 LPA", Year3: "₹6-10 LPA",
			Year4: "₹7-12 LPA", Year5: "₹8-15 LPA",
		},
		CareerStartup: {
			Year1: "₹0-5 LPA", Year2: "₹0-10 LPA", Year3: "₹5-20 LPA",
			Year4: "₹10-40 LPA", Year5: "₹15-100+ LPA",
		},
		CareerHigherStudies: {
			Year1: "₹0 (Stipend ₹30-50K/mo)", Year2: "₹0 (Stipend ₹35-60K/mo)", Year3: "₹8-15 LPA",
			Year4: "₹10-20 LPA", Year5: "₹14-30 LPA",
		},
		CareerMSAbroad: {
			Year1: "$0 (Studying)", Year2: "$60-90K/year", Year3: "$75-120K/year",
			Year4: "$90-150K/year", Year5: "$100-180K/year",
		},
	}

	if p, ok := projections[career]; ok {
		return p
	}
	return projections[CareerIT]
}

// GetRoadmap returns the preparation roadmap for a career.
func GetRoadmap(career Career) []dto.RoadmapStep {
	roadmaps := map[Career][]dto.RoadmapStep{
		CareerIT: {
			{Step: 1, Title: "Learn Programming Fundamentals", Description: "Master one language (Python/Java/JavaScript). Complete DSA basics on LeetCode/GeeksForGeeks.", Duration: "3 months"},
			{Step: 2, Title: "Build Projects & Portfolio", Description: "Build 3-5 real projects. Create GitHub portfolio. Learn Git, APIs, databases.", Duration: "3 months"},
			{Step: 3, Title: "Learn Frameworks & Tools", Description: "Pick a stack (MERN/Spring Boot/Django). Learn Docker, cloud basics (AWS/GCP).", Duration: "2 months"},
			{Step: 4, Title: "DSA & Interview Prep", Description: "Solve 200+ LeetCode problems. Practice system design. Mock interviews.", Duration: "3 months"},
			{Step: 5, Title: "Apply & Network", Description: "Apply on LinkedIn, Naukri, AngelList. Attend hackathons. Get referrals.", Duration: "1 month"},
		},
		CareerMBA: {
			{Step: 1, Title: "CAT/XAT/GMAT Preparation", Description: "Join coaching (IMS/TIME/CL) or self-study. Target 95+ percentile in CAT.", Duration: "6-8 months"},
			{Step: 2, Title: "Build Profile", Description: "Gain 2-3 years work experience. Get leadership roles. Volunteer work.", Duration: "Ongoing"},
			{Step: 3, Title: "Application & Essays", Description: "Research IIMs, XLRI, ISB, FMS. Write compelling SOPs and essays.", Duration: "2 months"},
			{Step: 4, Title: "GD/PI Preparation", Description: "Current affairs, case studies, mock GDs and PIs.", Duration: "2 months"},
			{Step: 5, Title: "Specialization Planning", Description: "Research Finance, Marketing, Operations, HR tracks. Network with alumni.", Duration: "1 month"},
		},
		CareerGovt: {
			{Step: 1, Title: "Choose Your Exam", Description: "UPSC CSE, SSC CGL, Banking (IBPS/SBI), State PSC, Railways. Pick based on your eligibility.", Duration: "1 month"},
			{Step: 2, Title: "Foundation Building", Description: "NCERT books (6-12), basic GK, aptitude. Join coaching if needed (Unacademy/BYJU's).", Duration: "3 months"},
			{Step: 3, Title: "Subject Deep Dive", Description: "Cover full syllabus. Make notes. Previous year papers analysis.", Duration: "6 months"},
			{Step: 4, Title: "Test Series & Revision", Description: "Join test series. Weekly full-length mocks. Analyze mistakes.", Duration: "3 months"},
			{Step: 5, Title: "Prelims → Mains → Interview", Description: "Clear each stage. Personality test prep for UPSC. Document verification.", Duration: "6-12 months"},
		},
		CareerStartup: {
			{Step: 1, Title: "Ideation & Validation", Description: "Identify problems worth solving. Talk to 50+ potential customers. Validate demand.", Duration: "2 months"},
			{Step: 2, Title: "MVP Development", Description: "Build minimum viable product. Use no-code tools if needed. Get first 10 users.", Duration: "2 months"},
			{Step: 3, Title: "Early Traction", Description: "Get to 100+ users. Iterate based on feedback. Find product-market fit.", Duration: "3 months"},
			{Step: 4, Title: "Funding & Team", Description: "Apply to incubators (IIT, NSRCEL, T-Hub). Pitch to angels. Build core team.", Duration: "3 months"},
			{Step: 5, Title: "Scale & Growth", Description: "Optimize unit economics. Hiring. Series A preparation. Scale marketing.", Duration: "6 months"},
		},
		CareerHigherStudies: {
			{Step: 1, Title: "Choose Exam & Specialization", Description: "GATE, NET, JAM, or direct admission. Pick M.Tech/M.Sc/PhD path.", Duration: "1 month"},
			{Step: 2, Title: "Exam Preparation", Description: "GATE: Focus on core subjects + aptitude. Target AIR under 500 for IITs.", Duration: "6 months"},
			{Step: 3, Title: "College Selection", Description: "Research IITs, IISc, NITs, IIITs. Check placement records and research labs.", Duration: "1 month"},
			{Step: 4, Title: "Research & Thesis", Description: "Choose research area. Publish papers. Build academic network.", Duration: "12-18 months"},
			{Step: 5, Title: "Placement/PhD Application", Description: "Campus placements or apply for PhD positions. Build research profile.", Duration: "3 months"},
		},
		CareerMSAbroad: {
			{Step: 1, Title: "GRE & TOEFL/IELTS Prep", Description: "Target GRE 320+, TOEFL 100+ or IELTS 7.5+. Use Magoosh/ETS material.", Duration: "3 months"},
			{Step: 2, Title: "University Shortlisting", Description: "Research universities (US/Canada/Germany/UK). Check admit chances on Yocket/Admits.fyi.", Duration: "2 months"},
			{Step: 3, Title: "SOP, LORs & Application", Description: "Write compelling SOPs. Get 3 strong LORs. Apply to 8-12 universities.", Duration: "3 months"},
			{Step: 4, Title: "Funding & Visa", Description: "Apply for scholarships, TA/RA positions. Education loan. F1/student visa.", Duration: "3 months"},
			{Step: 5, Title: "Pre-Departure", Description: "Housing, bank account, health insurance. Connect with seniors at target university.", Duration: "2 months"},
		},
	}

	if r, ok := roadmaps[career]; ok {
		return r
	}
	return roadmaps[CareerIT]
}

// GetRequiredSkills returns the skills needed for a career.
func GetRequiredSkills(career Career) []string {
	skills := map[Career][]string{
		CareerIT: {
			"Data Structures & Algorithms", "Programming (Python/Java/JS)",
			"Web Development (React/Node)", "Database Management (SQL/NoSQL)",
			"System Design", "Cloud Computing (AWS/GCP)",
			"Version Control (Git)", "Problem Solving",
		},
		CareerMBA: {
			"Quantitative Aptitude", "Verbal Ability & Reading Comprehension",
			"Logical Reasoning", "Data Interpretation",
			"Leadership & Teamwork", "Communication Skills",
			"Business Acumen", "Current Affairs",
		},
		CareerGovt: {
			"General Knowledge & Current Affairs", "Quantitative Aptitude",
			"English Language", "Logical Reasoning",
			"Indian Polity & Constitution", "Indian Economy",
			"History & Geography", "Essay Writing",
		},
		CareerStartup: {
			"Product Thinking", "Sales & Marketing",
			"Financial Planning", "Leadership & Team Building",
			"Technical Skills (Full-Stack/No-Code)", "Fundraising & Pitching",
			"Customer Development", "Growth Hacking",
		},
		CareerHigherStudies: {
			"Core Subject Expertise", "Research Methodology",
			"Academic Writing", "GATE/NET Exam Skills",
			"Programming (for CS/IT)", "Lab Work & Experimentation",
			"Paper Reading & Review", "Presentation Skills",
		},
		CareerMSAbroad: {
			"GRE Verbal & Quant", "TOEFL/IELTS English Proficiency",
			"Research Experience", "Academic Writing (SOP)",
			"Core Domain Knowledge", "Programming & Tools",
			"Networking & Communication", "Cross-Cultural Adaptability",
		},
	}

	if s, ok := skills[career]; ok {
		return s
	}
	return skills[CareerIT]
}

// GetSuggestedExams returns relevant exams for a career.
func GetSuggestedExams(career Career) []string {
	exams := map[Career][]string{
		CareerIT:            {"GATE CS", "Google Kickstart", "CodeChef/Codeforces", "AWS Certification", "Company-specific OAs"},
		CareerMBA:           {"CAT", "XAT", "GMAT", "NMAT", "SNAP", "IIFT"},
		CareerGovt:          {"UPSC CSE", "SSC CGL", "IBPS PO", "SBI PO", "RBI Grade B", "State PSC"},
		CareerStartup:       {"No specific exams - focus on building", "Y Combinator Application", "Shark Tank India (if applicable)"},
		CareerHigherStudies: {"GATE", "UGC NET", "CSIR NET", "IIT JAM", "JEST"},
		CareerMSAbroad:      {"GRE General", "TOEFL iBT", "IELTS Academic", "GRE Subject (optional)"},
	}

	if e, ok := exams[career]; ok {
		return e
	}
	return exams[CareerIT]
}

// GetSuggestedColleges returns recommended institutions for a career.
func GetSuggestedColleges(career Career) []string {
	colleges := map[Career][]string{
		CareerIT: {
			"IIT Bombay/Delhi/Madras (B.Tech/M.Tech)",
			"NIT Trichy/Warangal/Surathkal",
			"IIIT Hyderabad / BITS Pilani",
			"Top product companies (Google, Microsoft, Amazon)",
		},
		CareerMBA: {
			"IIM Ahmedabad / Bangalore / Calcutta",
			"IIM Lucknow / Indore / Kozhikode",
			"XLRI Jamshedpur / FMS Delhi",
			"ISB Hyderabad / IIM Udaipur (1-year)",
		},
		CareerGovt: {
			"LBSNAA (IAS Training)", "SVPNPA (IPS Training)",
			"Reserve Bank of India", "State Administrative Services",
		},
		CareerStartup: {
			"IIT/IIM Incubators", "NSRCEL (IIM Bangalore)",
			"T-Hub Hyderabad", "Startup India Hub",
			"Y Combinator / Techstars (global)",
		},
		CareerHigherStudies: {
			"IISc Bangalore", "IIT Bombay/Delhi/Madras/Kanpur",
			"TIFR / ISI Kolkata", "JNU / Delhi University",
		},
		CareerMSAbroad: {
			"MIT / Stanford / CMU (US)",
			"UC Berkeley / Georgia Tech / UIUC",
			"ETH Zurich / TU Munich (Europe)",
			"University of Toronto / UBC (Canada)",
		},
	}

	if c, ok := colleges[career]; ok {
		return c
	}
	return colleges[CareerIT]
}
