// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package rdp

type SectionId struct {
	UnitId string
}

// AcceptSectionId returns a section id if possible
func (p *Parser) AcceptSectionId() *SectionId {
	junk, _ := AcceptJunk(p.input[p.pos:])
	p.pos += len(junk)
	if p.iseof() {
		return nil
	}
	// if we are here, we have a section id.
	if match := rxCourierHeader.FindSubmatch(p.input[p.pos:]); match != nil {
		// we have a courier header. capture it and advance the position.
		p.pos += len(match[0])
		return &SectionId{
			UnitId: string(match[1]),
		}
	} else if match := rxElementHeader.FindSubmatch(p.input[p.pos:]); match != nil {
		// we have an element header. capture it and advance the position.
		p.pos += len(match[0])
		return &SectionId{
			UnitId: string(match[1]),
		}
	} else if match := rxFleetHeader.FindSubmatch(p.input[p.pos:]); match != nil {
		// we have a fleet header. capture it and advance the position.
		p.pos += len(match[0])
		return &SectionId{
			UnitId: string(match[1]),
		}
	} else if match := rxGarrisonHeader.FindSubmatch(p.input[p.pos:]); match != nil {
		// we have a garrison header. capture it and advance the position.
		p.pos += len(match[0])
		return &SectionId{
			UnitId: string(match[1]),
		}
	} else if match := rxTribeHeader.FindSubmatch(p.input[p.pos:]); match != nil {
		// we have a tribe header. capture it and advance the position.
		p.pos += len(match[0])
		return &SectionId{
			UnitId: string(match[1]),
		}
	}
	// we didn't find a section id, so panic
	panic("section id failed")
}
