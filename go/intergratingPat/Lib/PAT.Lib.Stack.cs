using System;
using System.Collections.Generic;
using System.Text;
using PAT.Common.Classes.Expressions.ExpressionClass;

//the namespace must be PAT.Lib, the class and method names can be arbitrary
namespace PAT.Lib
{
    public class Stack: ExpressionValue
    {

        public System.Collections.Generic.Stack<int> stack;

        //default constructor
        public Stack()
        {
            this.stack = new System.Collections.Generic.Stack<int>();
        }

        public Stack(System.Collections.Generic.Stack<int> stack)
        {
            this.stack = stack;
        }
        
        //override
        public override string ToString()
        {
            return "[" + ExpressionID + "]";
        }

        public override ExpressionValue GetClone()
        {
            return new Stack(new System.Collections.Generic.Stack<int>(stack));
        }


        public override string ExpressionID
        {
            get
            {
                String returnString = "";
                foreach (int element in stack)
                {
                    returnString += element.ToString() + ",";
                }

                if (returnString.Length > 0)
                {
                    returnString = returnString.Substring(0, returnString.Length - 1);
                }

                return returnString;
            }
        }

        //public override string GetID()
        //{
        //    String returnString = "";
        //    foreach (int element in stack) 
        //    {
        //        returnString += element.ToString() + ",";
        //    }

        //    if (returnString.Length > 0)
        //    {
        //        returnString = returnString.Substring(0, returnString.Length - 1);
        //    }

        //    return returnString;
        //}

        public void Push(int element)
        {
            this.stack.Push(element);
        }

        public int Pop()
        {
            if (this.stack.Count > 0)
            {
                return stack.Pop();
            }
            else
            {
                
                //throw PAT Runtime exception
                throw new RuntimeException("Access an empty stack!");
            }

        }

        public bool Contains(int element)
        {
            return this.stack.Contains(element);
        }

        public int Peek()
        {
            if (stack.Count > 0)
            {
                return stack.Peek();
            }
            else
            {
                //throw PAT Runtime exception
                throw new RuntimeException("Access an empty stack!");
            }
        }

        public void Clear()
        {
            this.stack.Clear();
        }

        public int Count()
        {
            return this.stack.Count;
        }

    }
}
