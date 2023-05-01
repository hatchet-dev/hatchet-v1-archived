import Spinner from "../loaders";
import React from "react";
import { StyledPlaceholder } from "./styles";

export type Props = {
  width?: string;
  height?: string;
  background?: string;
  loading?: boolean;
  children?: React.ReactNode;
};

const Placeholder: React.FC<Props> = ({
  width,
  height,
  background,
  children,
  loading,
}) => {
  return (
    <StyledPlaceholder width={width} height={height} background={background}>
      {loading ? <Spinner /> : children}
    </StyledPlaceholder>
  );
};

export default Placeholder;
