import SectionArea from "../sectionarea";
import Logo from "../logo";
import React from "react";

import { LogoContainer } from "./styles";

export type Props = {
  width?: string;
  height?: string;
  background?: string;
  loading?: boolean;
  children?: React.ReactNode;
};

const SectionAreaWithLogo: React.FC<Props> = ({
  width,
  height,
  background,
  children,
  loading,
}) => {
  return (
    <SectionArea
      width={width}
      height={height}
      background={background}
      loading={loading}
    >
      <LogoContainer>
        <Logo />
      </LogoContainer>
      {children}
    </SectionArea>
  );
};

export default SectionAreaWithLogo;
