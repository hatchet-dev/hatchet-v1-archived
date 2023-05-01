import React from "react";
import { CenteredLogo, CenteredLogoContainer } from "./styles";
import hatchet from "../../assets/hatchet.png";

export type Props = {
  width?: string;
  height?: string;
  padding?: string;
};

const Logo: React.FC<Props> = ({ width, height, padding }) => {
  return (
    <CenteredLogoContainer width={width} height={height} padding={padding}>
      <CenteredLogo src={hatchet} />
    </CenteredLogoContainer>
  );
};

export default Logo;
