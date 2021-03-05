const calculator = {
  displayValue: '0',
  firstOperand: null,
  waitingForSecondOperand: false,
  operator: null,
};

function inputDigit(digit) {
  const { displayValue, waitingForSecondOperand } = calculator;

  if (waitingForSecondOperand === true) {
    calculator.displayValue = digit;
    
    calculator.waitingForSecondOperand = false;
  } else {
    calculator.displayValue = displayValue === '0' ? digit : displayValue + digit;
  }
}

function inputDecimal(dot) {
  // If the `displayValue` does not contain a decimal point
  if (!calculator.displayValue.includes(dot)) {
    // Append the decimal point
    calculator.displayValue += dot;
  }
}

function handleOperator(nextOperator) {
  const { firstOperand, displayValue, operator } = calculator
  const inputValue = parseFloat(displayValue);
  
  if (operator && calculator.waitingForSecondOperand)  {
    calculator.operator = nextOperator;
    return;
  }

  if (firstOperand == null) {
    calculator.firstOperand = inputValue;
  } else if (operator) {
      const currentValue = firstOperand || 0;
      const result = Processing(operator,currentValue,inputValue)
      result.then((data)=>{
        document.getElementById("displayresult").value = data;
        calculator.displayValue = String(data);
        calculator.firstOperand = data;
      })
  }

  calculator.waitingForSecondOperand = true;
  calculator.operator = nextOperator;
}

async function Processing(operator,currentValue,inputValue){
  const settings = {
    method: 'POST',
    body: JSON.stringify({
        "operator": operator,
        "input1": currentValue.toString(),
        "input2":inputValue.toString(),
        "result":"",
    }),
    headers: {
         Accept: 'application/json',
        'Content-Type': 'application/json',
    } 
  };

  const response = await fetch('/calculator',settings)
  const data = await response.json();
  return data.result;
}

function resetCalculator() {
  calculator.displayValue = '0';
  calculator.firstOperand = null;
  calculator.waitingForSecondOperand = false;
  calculator.operator = null;
}

function updateDisplay() {
  const display = document.querySelector('.calculator-screen');
  display.value = calculator.displayValue;
}

updateDisplay();

const keys = document.querySelector('.calculator-keys');
keys.addEventListener('click', (event) => {
  const { target } = event;
  if (!target.matches('button')) {
    return;
  }


  if (target.classList.contains('operator')) {

    handleOperator(target.value);
		updateDisplay();
    return;
  }

  if (target.classList.contains('decimal')) {
    inputDecimal(target.value);
		updateDisplay();
    return;
  }

  if (target.classList.contains('all-clear')) {
    resetCalculator();
		updateDisplay();
    return;
  }

  inputDigit(target.value);
  updateDisplay();
});