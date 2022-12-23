import { Relative, Span } from "components/globals";
import React, { useEffect, useRef, useState } from "react";
import { StyledSectionArea } from "./styles";

export type Props = {
  width?: number;
  height?: number;
  background?: string;
};

const SectionArea: React.FC<Props> = ({
  width,
  height,
  background,
  children,
}) => {
  return (
    <StyledSectionArea width={width} height={height} background={background}>
      {children}
    </StyledSectionArea>
  );
};

export default SectionArea;
