import React, { Component, useRef } from "react";
import styled from "styled-components";
import AceEditor from "react-ace";

import "./ace-hatchet-theme.js";
import "ace-builds/src-noconflict/mode-json";

type Props = {
  value: string;
  onChange?: (e: any) => void;
  readOnly?: boolean;
  height?: string;
};

const CodeBlock: React.FC<Props> = ({ value, onChange, readOnly, height }) => {
  const handleChange = (e: any) => {
    onChange && onChange(e);
  };

  return (
    <AceContainer>
      <AceEditor
        mode="json"
        value={value}
        theme="hatchet"
        onChange={handleChange}
        name="codeEditor"
        readOnly={readOnly}
        editorProps={{ $blockScrolling: true }}
        height={height}
        width="100%"
        style={{ borderRadius: "10px" }}
        showPrintMargin={false}
        showGutter={false}
        highlightActiveLine={true}
        fontSize={14}
      />
    </AceContainer>
  );
};

export default CodeBlock;

const AceContainer = styled.div`
  .ace_scrollbar {
    display: none;
  }
  .ace_editor,
  .ace_editor * {
    font-family: monospace;
    font-size: 12px !important;
    font-weight: 400 !important;
    letter-spacing: 0 !important;
  }

  height: 100%;
  width: 100%;
`;
