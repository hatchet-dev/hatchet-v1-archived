import StandardButton from "../buttons";
import CopyToClipboard from "../copytoclipboard";
import { FlexRow } from "../globals";
import React, { useState } from "react";
import { CodeBlock } from "./styles";

export type Props = {
  code_line: string;
};

const CopyCodeLine: React.FC<Props> = ({ code_line }) => {
  const [successCopy, setSuccessCopy] = useState(false);

  return (
    <FlexRow>
      <CodeBlock>{code_line}</CodeBlock>
      <CopyToClipboard
        as={StandardButton}
        wrapperProps={{
          label: successCopy ? "Copied!" : "Copy",
          style_kind: "muted",
        }}
        text={code_line}
        onSuccess={() => {
          setSuccessCopy(true);
        }}
      />
    </FlexRow>
  );
};

export default CopyCodeLine;
