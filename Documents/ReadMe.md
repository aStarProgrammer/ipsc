# IPSC Help

IPSC is a tool to create static html site

Description:

Commands:
* Create New Empty Site
	Run the following command to create empty site
		IPSC -Command NewSite -SiteTitle  -SiteAuthor  -SiteDescription  -OutputFolder -SiteFolder
	
	For Example:
		IPSC -Command "NewSite" -SiteTitle "Test Site" -SiteAuthor "Chao(sdxianchao@gmail.com)" -SiteDescription "Test Site for IPSC" -OutputFolder "F:\SiteOutputFolder" -SiteFolder "F:\TestSite"
	
	Above command will create an empty site folder at F:\TestSite, and it looks like
		F:\TestSite
			○ Test Site.sp
			○ Src
				§ Markdown
				§ Html
			○ Output(Soft Link->F:\SiteOutputFolder)
				§ Pages
	
	Test Site.sp looks like
		{
		"Author": "Chao(sdxianchao@gmail.com)",
		"CreateTime": "2019-11-25 01:45:41",
		"Description": "Test Site for IPSC",
		"ID": "7f297e3cee64da063184fe0edd6d928b",
		"IndexPageSourceFile": {
			"Author": "",
			"CreateTime": "",
			"Description": "",
			"ID": "",
			"IsTop": false,
			"LastComplied": "",
			"LastModified": "",
			"OutputFile": 0,
			"SourceFilePath": "",
			"Status": "",
			"Title": "",
			"TitleImage": "",
			"Type": ""
		},
		"LastComplieSummary": "",
		"LastComplied": "",
		"LastModified": "2019-11-25 01:45:41",
		"MorePageSourceFiles": null,
		"OutputFiles": null,
		"OutputFolderPath": "F:\\SiteOutputFolder",
		"SourceFiles": null,
		"Title": "Test Site"
		}

