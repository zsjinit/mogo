import dataSourceMenuStyles from "@/pages/DataLogs/components/DataSourceMenu/index.less";
import SourceHeader from "@/pages/DataLogs/components/DataSourceMenu/SourceHeader";
import LoggingLibrary from "@/pages/DataLogs/components/DataSourceMenu/LoggingLibrary";
import { useEffect } from "react";
import { useModel } from "@@/plugin-model/useModel";
import classNames from "classnames";

const DataSourceMenu = () => {
  const { doGetDatabaseList } = useModel("dataLogs");
  const { foldingState } = useModel("dataLogs");
  useEffect(() => {
    doGetDatabaseList();
  }, []);

  return (
    <div
      className={classNames(
        dataSourceMenuStyles.dataSourceMenuMain,
        foldingState && dataSourceMenuStyles.dataSourceMenuHidden
      )}
    >
      <SourceHeader />
      <LoggingLibrary />
    </div>
  );
};

export default DataSourceMenu;
