import { FlexColCenter } from "../globals";
import Spinner from "../loaders";
import React from "react";
import { StyledSectionArea } from "./styles";

export type Props = {
  width?: string;
  height?: string;
  flex?: string;
  background?: string;
  border?: string;
  loading?: boolean;
  children?: React.ReactNode;
};

const SectionArea: React.FC<Props> = ({
  width,
  height,
  background,
  border,
  flex,
  children,
  loading,
}) => {
  return (
    <StyledSectionArea
      width={width}
      height={height}
      background={background}
      border={border}
      flex={flex}
    >
      {loading ? (
        <FlexColCenter>
          <Spinner />
        </FlexColCenter>
      ) : (
        children
      )}
    </StyledSectionArea>
  );
};

export default SectionArea;
