function Button() {
  function handleClick() {
    alert("good job");
  }
  return <button onClick={handleClick}>Submit</button>;
}

export default Button;
