interface dropDownSettings{
    value: string,
    callBack: (event: string) => void,
    optionObject: { type: string, name: string }[]
}

 export default function DropDownButton({ value, callBack, optionObject }: dropDownSettings) {
    return (
      <select
        className="inference"
        value={value}
        onChange={(event) => {
          callBack(event.target.value);
        }}
      >
        {optionObject.map((type) => {
          return (
            <option value={type.type} key={type.type}>
              {type.name}
            </option>
          );
        })}
      </select>
    );
  }