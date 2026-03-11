package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateTrackingDocx creates a minimal DOCX file with an embedded remote
// image that points to the tracking URL. When Microsoft Word opens the file
// and loads external images, it makes an HTTP GET to the tracking URL,
// registering an "Opened Attachment" event in CyberPhish.
//
// trackURL should be: https://your-domain.com/track-open?rid=RECIPIENT_RID
// outputPath is the destination file, e.g. ./static/attachments/Invoice.docx
func GenerateTrackingDocx(trackURL, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`,

		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`,

		"word/_rels/document.xml.rels": fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s" TargetMode="External"/>
</Relationships>`, trackURL),

		"word/document.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"
            xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"
            xmlns:o="urn:schemas-microsoft-com:office:office"
            xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
            xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math"
            xmlns:v="urn:schemas-microsoft-com:vml"
            xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
            xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
            xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
            xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture"
            mc:Ignorable="w14 wp14">
  <w:body>
    <w:p>
      <w:r>
        <w:drawing>
          <wp:inline distT="0" distB="0" distL="0" distR="0">
            <wp:extent cx="1" cy="1"/>
            <wp:docPr id="1" name="tracker"/>
            <a:graphic>
              <a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">
                <pic:pic>
                  <pic:nvPicPr>
                    <pic:cNvPr id="0" name="tracker"/>
                    <pic:cNvPicPr/>
                  </pic:nvPicPr>
                  <pic:blipFill>
                    <a:blip r:link="rId1"/>
                    <a:stretch><a:fillRect/></a:stretch>
                  </pic:blipFill>
                  <pic:spPr>
                    <a:xfrm><a:off x="0" y="0"/><a:ext cx="1" cy="1"/></a:xfrm>
                    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
                  </pic:spPr>
                </pic:pic>
              </a:graphicData>
            </a:graphic>
          </wp:inline>
        </w:drawing>
      </w:r>
    </w:p>
    <w:sectPr/>
  </w:body>
</w:document>`,
	}

	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			return fmt.Errorf("creating zip entry %s: %w", name, err)
		}
		if _, err = f.Write([]byte(content)); err != nil {
			return fmt.Errorf("writing zip entry %s: %w", name, err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("closing zip: %w", err)
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

// GenerateTrackingXlsx creates a minimal XLSX file with an embedded external
// image reference that triggers the tracking URL when opened in Excel.
//
// trackURL should be: https://your-domain.com/track-open?rid=RECIPIENT_RID
func GenerateTrackingXlsx(trackURL, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
  <Override PartName="/xl/drawings/drawing1.xml" ContentType="application/vnd.openxmlformats-officedocument.drawing+xml"/>
</Types>`,

		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`,

		"xl/workbook.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"
          xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheets>
    <sheet name="Sheet1" sheetId="1" r:id="rId1"/>
  </sheets>
</workbook>`,

		"xl/_rels/workbook.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
</Relationships>`,

		"xl/worksheets/sheet1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"
           xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheetData/>
  <drawing r:id="rId1"/>
</worksheet>`,

		"xl/worksheets/_rels/sheet1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/drawing" Target="../drawings/drawing1.xml"/>
</Relationships>`,

		"xl/drawings/drawing1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<xdr:wsDr xmlns:xdr="http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing"
          xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
          xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
          xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
  <xdr:absoluteAnchor>
    <xdr:pos x="0" y="0"/>
    <xdr:ext cx="1" cy="1"/>
    <xdr:pic>
      <xdr:nvPicPr>
        <xdr:cNvPr id="2" name="tracker"/>
        <xdr:cNvPicPr/>
      </xdr:nvPicPr>
      <xdr:blipFill>
        <a:blip r:link="rId1"/>
        <a:stretch><a:fillRect/></a:stretch>
      </xdr:blipFill>
      <xdr:spPr>
        <a:xfrm><a:off x="0" y="0"/><a:ext cx="1" cy="1"/></a:xfrm>
        <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
      </xdr:spPr>
    </xdr:pic>
    <xdr:clientData/>
  </xdr:absoluteAnchor>
</xdr:wsDr>`,

		"xl/drawings/_rels/drawing1.xml.rels": fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s" TargetMode="External"/>
</Relationships>`, trackURL),
	}

	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			return fmt.Errorf("creating zip entry %s: %w", name, err)
		}
		if _, err = f.Write([]byte(content)); err != nil {
			return fmt.Errorf("writing zip entry %s: %w", name, err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("closing zip: %w", err)
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}
