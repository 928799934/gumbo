package gumbo

/*
#cgo linux LDFLAGS: -lgumbo
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <gumbo.h>

void test(){
	GumboOutput* output = gumbo_parse("<body>Hello, World!</body>");
	GumboNode* root = output->root;
	int i = 0;
	GumboVector *children = &root->v.element.children;
	for ( ; i < children->length; i++ ) {
		root = (GumboNode *)children->data[i];
		if (root->v.element.tag != GUMBO_TAG_BODY) {
			continue;
		}
		GumboVector *child = &root->v.element.children;
		root = (GumboNode *)child->data[0];
		GumboText *t = &root->v.text;
		printf("%s\n",t->text);
	}
	gumbo_destroy_output(&kGumboDefaultOptions, output);
}

GumboOutput*
parse(
	const char* html,
	size_t len
) {
	return gumbo_parse_with_options(&kGumboDefaultOptions,html,len);
}

void
destroy(
	GumboOutput* output
) {
	gumbo_destroy_output(&kGumboDefaultOptions, output);
}

bool
getNodeByTag(
	GumboNode* root,
	GumboTag tag,
	GumboNode** node
) {
	if (root->type != GUMBO_NODE_ELEMENT)
		return false;

	if (root->v.element.tag == tag) {
		*node = root;
		return true;
	}

	unsigned int i = 0;
	GumboVector *children = &root->v.element.children;
	for ( ; i < children->length; i++ ) {
		if( getNodeByTag(
			(GumboNode *)children->data[i],
			tag,
			node)
		) {
			return true;
		}
	}
	return false;
}

bool
getNodeByTagAndAttr(
	GumboNode* root,
	GumboTag tag,
	const char* name,
	const char* value,
	GumboNode** node
) {
	if (root->type != GUMBO_NODE_ELEMENT)
		return NULL;

	GumboVector *children = &root->v.element.children;
	if (root->v.element.tag == tag) {

		GumboAttribute *attr = gumbo_get_attribute(
			&root->v.element.attributes,
			name);
		if( attr != NULL &&
			memcmp(value,attr->value,strlen(value)) == 0
		) {
			*node = root;
			return true;
		}

		if(children->length == 0)
			return false;
	}

	unsigned int i = 0;
	for ( ; i < children->length; i++ ) {
		if( getNodeByTagAndAttr(
			(GumboNode *)children->data[i],
			tag,
			name,
			value,
			node
			)
		) {
			return true;
		}
	}
	return false;
}

bool
getText(
	GumboNode* root,
	const char** value
) {
	if (root->type != GUMBO_NODE_ELEMENT)
		return false;

	GumboVector *child = &root->v.element.children;

	if (child->length != 1)
		return false;

	root = (GumboNode *)child->data[0];
	GumboText *t = &root->v.text;
	*value = t->text;
	return true;
}

bool
getAttribute(
	GumboNode* root,
	const char* attr,
	const char** value
) {
	if (root->type != GUMBO_NODE_ELEMENT)
		return false;
	GumboVector* attrs = &root->v.element.attributes;
	GumboAttribute *attribute = gumbo_get_attribute(attrs,attr);
	if (attribute != NULL) {
		*value = attribute->value;
		return true;
	}
	return false;
}

*/
import "C"

import "unsafe"

type Gumbo struct {
	output *C.GumboOutput
}

func NewGumbo() *Gumbo {
	return &Gumbo{}
}

func NewGumboParse(html string) *Gumbo {
	gb := &Gumbo{}
	gb.Parse(html)
	return gb
}

func (this *Gumbo) Parse(html string) {
	cs := C.CString(html)
	defer C.free(unsafe.Pointer(cs))

	if this.output != nil {
		this.Destory()
	}

	this.output = C.parse(cs, C.size_t(len(html)))
}

func (this *Gumbo) Destory() {
	if this.output != nil {
		C.destroy(this.output)
		this.output = nil
	}
}

func (this *Gumbo) GetNodeByTag(tag C.GumboTag, root *C.GumboNode) (*C.GumboNode, bool) {
	var node *C.GumboNode

	if root == nil {
		root = this.output.root
	}

	ok := C.getNodeByTag(root, tag, &node)
	return node, bool(ok)
}

func (this *Gumbo) GetNodeByTagAndAttr(tag C.GumboTag, attrName, attrValue string, root *C.GumboNode) (*C.GumboNode, bool) {
	szAttrName := C.CString(attrName)
	defer C.free(unsafe.Pointer(szAttrName))
	szAttrValue := C.CString(attrValue)
	defer C.free(unsafe.Pointer(szAttrValue))

	var node *C.GumboNode

	if root == nil {
		root = this.output.root
	}

	ok := C.getNodeByTagAndAttr(root, tag, szAttrName, szAttrValue, &node)
	return node, bool(ok)
}

func (this *Gumbo) GetText(node *C.GumboNode) (string, bool) {
	szText := C.CString("")
	defer C.free(unsafe.Pointer(szText))

	ok := C.getText(node, &szText)
	return C.GoString(szText), bool(ok)
}

func (this *Gumbo) GetAttribute(node *C.GumboNode, attrName string) (string, bool) {
	szAttrValue := C.CString("")
	defer C.free(unsafe.Pointer(szAttrValue))
	szAttrName := C.CString(attrName)
	defer C.free(unsafe.Pointer(szAttrName))

	ok := C.getAttribute(node, szAttrName, &szAttrValue)
	return C.GoString(szAttrValue), bool(ok)
}
