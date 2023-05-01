import {
  HorizontalSpacer,
  P,
  FlexColScroll,
  ErrorBar,
} from "hatchet-components";
import React, { useState, useEffect } from "react";
import CodeBlock from "components/codeblock";

type Props = {
  current_values?: string;
  set_values: (vals: string) => void;
  jsonParseErr?: string;
};

const UploadJSONValues: React.FC<Props> = ({
  current_values,
  set_values,
  jsonParseErr,
}) => {
  const [vals, setVals] = useState(current_values);

  useEffect(() => {
    set_values(vals);
  }, [vals]);

  return (
    <>
      <HorizontalSpacer spacepixels={24} />
      <P>Upload your JSON variables here.</P>
      <HorizontalSpacer spacepixels={24} />
      <FlexColScroll height="200px" width="100%">
        <CodeBlock
          value={vals}
          height="200px"
          onChange={(e) => {
            setVals(e);
          }}
        />
      </FlexColScroll>
      {jsonParseErr && <HorizontalSpacer spacepixels={20} />}
      {jsonParseErr && <ErrorBar text={jsonParseErr} />}
    </>
  );
};

export default UploadJSONValues;
