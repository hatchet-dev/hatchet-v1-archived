import { FlexColCenter } from "components/globals";
import Spinner from "components/loaders";
import React from "react";
import { StyledSectionArea } from "./styles";

export type Props = {
  width?: number;
  height?: number;
  background?: string;
  loading?: boolean;
};

const SectionArea: React.FC<Props> = ({
  width,
  height,
  background,
  children,
  loading,
}) => {
  return (
    <StyledSectionArea
      width={width && width + "px"}
      height={height && height + "px"}
      background={background}
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
