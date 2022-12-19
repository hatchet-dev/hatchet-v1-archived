import { Relative, Span } from "components/globals";
import React, { useEffect, useRef, useState } from "react";
import { StyledSectionArea } from "./styles";

export type Props = {
  width?: number;
};

const SectionArea: React.FC<Props> = ({ width, children }) => {
  return <StyledSectionArea width={width}>{children}</StyledSectionArea>;
};

export default SectionArea;
