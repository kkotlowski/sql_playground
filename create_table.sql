USE [PanoramaFirm]
GO

/****** Object:  Table [dbo].[Companies]    Script Date: 11/20/2021 8:31:57 AM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[Companies](
	[ID] [int] IDENTITY(1,1) NOT NULL,
	[MainBranch] [varchar](50) NULL,
	[Branch] [varchar](50) NULL,
	[Name] [varchar](50) NULL,
	[Website] [varchar](50) NULL,
	[Mail] [varchar](50) NULL,
	[PhoneNumber] [varchar](50) NULL
) ON [PRIMARY]
GO


