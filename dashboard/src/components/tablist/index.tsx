import React, { useState } from "react";
import { Tab, TabWrapper } from "./styles";

export type Props = {
  tabs: string[];
  selectTab: (tab: string) => void;
};

const TabList: React.FC<Props> = ({ tabs, selectTab }) => {
  const [selectedTab, setSelectedTab] = useState(tabs[0]);

  const setTab = (tab: string) => {
    setSelectedTab(tab);
    selectTab(tab);
  };

  return (
    <TabWrapper>
      {tabs.map((val, i) => {
        return (
          <Tab selected={val == selectedTab} onClick={() => setTab(val)}>
            {val}
          </Tab>
        );
      })}
    </TabWrapper>
  );
};

export default TabList;
