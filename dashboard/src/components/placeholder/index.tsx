import Spinner from "components/loaders";
import React from "react";
import { StyledPlaceholder } from "./styles";

export type Props = {
  width?: number;
  height?: number;
  background?: string;
  loading?: boolean;
};

const Placeholder: React.FC<Props> = ({
  width,
  height,
  background,
  children,
  loading,
}) => {
  return (
    <StyledPlaceholder
      width={width && width + "px"}
      height={height && height + "px"}
      background={background}
    >
      {loading ? <Spinner /> : children}
    </StyledPlaceholder>
  );
};

export default Placeholder;
