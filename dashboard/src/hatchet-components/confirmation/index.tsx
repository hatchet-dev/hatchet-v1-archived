import StandardButton from "../buttons";
import CopyToClipboard from "../copytoclipboard";
import {
  FlexCol,
  FlexRow,
  FlexRowRight,
  HorizontalSpacer,
  P,
} from "../globals";
import TextInput from "../textinput";
import React, { useState } from "react";

export type Props = {
  confirm_text: string;
  prompt: string;
  button_label: string;
  confirm_text_example?: string;
  confirmed: () => void;
  button_material_icon?: string;
};

const Confirmation: React.FC<Props> = ({
  confirm_text,
  prompt,
  confirm_text_example,
  button_label,
  button_material_icon = "delete",
  confirmed,
}) => {
  const [confirmText, setConfirmText] = useState("");

  return (
    <FlexCol>
      <P>{prompt}</P>
      <HorizontalSpacer spacepixels={20} />
      <TextInput
        placeholder={`ex. ${confirm_text_example || "delete"}`}
        type="text"
        width="100%"
        on_change={(val) => {
          setConfirmText(val);
        }}
      />
      <HorizontalSpacer spacepixels={20} />
      <FlexRowRight>
        <StandardButton
          label={button_label}
          material_icon={button_material_icon}
          icon_side="right"
          on_click={() => {
            confirmed();
          }}
          margin={"0"}
          disabled={confirm_text != confirmText}
        />
      </FlexRowRight>
    </FlexCol>
  );
};

export default Confirmation;
