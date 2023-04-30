import { MaterialIcon } from "../globals";
import Spinner from "../loaders";
import React from "react";
import { StyledDisplayButton, StyledStandardButton } from "./styles";

export type Props = {
  material_icon?: string;
  id?: string;
  label: string;
  size?: "small" | "default";
  icon_side?: "left" | "right";
  style_kind?: "default" | "muted";
  margin?: string;
  disabled?: boolean;
  on_click?: () => void;
  is_loading?: boolean;
};

const StandardButton: React.FC<Props> = ({
  id,
  material_icon,
  label,
  size = "default",
  style_kind = "default",
  icon_side = "left",
  margin = "0 8px",
  disabled = false,
  on_click,
  is_loading = false,
}) => {
  if (is_loading) {
    const children = [<Spinner />, label];

    return (
      <StyledStandardButton
        id={id}
        size={size}
        has_icon={true}
        onClick={on_click}
        icon_side={icon_side}
        margin={margin}
        disabled={disabled}
        style_kind={style_kind}
      >
        {icon_side == "left" ? children : children.reverse()}
      </StyledStandardButton>
    );
  }

  if (material_icon) {
    const children = [
      <MaterialIcon key="0" className="material-icons">
        {material_icon}
      </MaterialIcon>,
      label,
    ];

    return (
      <StyledStandardButton
        id={id}
        size={size}
        has_icon={true}
        onClick={on_click}
        icon_side={icon_side}
        margin={margin}
        disabled={disabled}
        style_kind={style_kind}
      >
        {icon_side == "left" ? children : children.reverse()}
      </StyledStandardButton>
    );
  }

  return (
    <StyledStandardButton
      id={id}
      size={size}
      has_icon={false}
      onClick={on_click}
      disabled={disabled}
      style_kind={style_kind}
    >
      {label}
    </StyledStandardButton>
  );
};

export default StandardButton;

export type DisplayButtonProps = {
  label: string;
  on_click?: () => void;
};

export const DisplayButton: React.FC<DisplayButtonProps> = ({
  label,
  on_click,
}) => {
  return <StyledDisplayButton onClick={on_click}>{label}</StyledDisplayButton>;
};
