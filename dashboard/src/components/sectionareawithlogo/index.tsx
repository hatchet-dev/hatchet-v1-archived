import SectionArea from "components/sectionarea";
import Logo from "components/logo";
import React from "react";

import { LogoContainer } from "./styles";

export type Props = {
  width?: number;
  height?: number;
  background?: string;
  loading?: boolean;
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
