import { useEffect, useRef } from "react";
import usePrevious from "./useprevious";
import isEqual from "lodash.isequal";
import cloneDeep from "lodash.clonedeep";

// useIsModified returns true if the value has been modified from it's first non-undefined
// value, false otherwise.
// note that we create deep clones of the values in case array refs are passed, in which case
// comparison would always return true
function useIsModified(value: any, originalValue?: any) {
  const ref = useRef();
  let isModified;

  function reset(reset?: any) {
    ref.current = cloneDeep(reset);
  }

  if (!originalValue && !value) {
    isModified = false;
  } else {
    if (!ref.current) {
      if (originalValue) {
        ref.current = cloneDeep(originalValue);
      } else {
        ref.current = cloneDeep(value);
      }
    }

    isModified = !isEqual(ref.current, value);
  }

  return { isModified, reset };
}

export default useIsModified;
