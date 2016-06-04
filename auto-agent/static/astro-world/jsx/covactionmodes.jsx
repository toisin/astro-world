/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var TwoRecordSelection = React.createClass({
  // getInitialState: function() {
  //   return {mode: 0};
  // },

  render: function() {
      var state = this.state;
      var user = this.props.user;
      var app = this.props.app;
      var prompt = user.getPrompt();

      return  <div className ="hbox"><div>
        <table>
          <tbody>
          <tr>
            <td>&nbsp;</td>
            <td colspan="3" className="question">First Record</td>
          </tr>
          <tr>
            <td>Fitness</td>
            <td><label><img src="graphics/excellent fitness.jpg"/><br/>
            <input type="radio" name="r1">Excellent</input></label></td>
            <td>&nbsp;</td>
            <td><label><img src="graphics/average fitness.jpg"/><br/>
            <input type="radio" name="r1">Average</input></label></td>
          </tr>
          <tr>
            <td>Parents health</td>
            <td><label><img src="graphics/excellent parents.jpg"/><br/>
            <input type="radio" name="r2">Excellent</input></label></td>
            <td>&nbsp;</td>
            <td><label><img src="graphics/fair parents.jpg"/><br/>
            <input type="radio" name="r2">Fair</input></label></td>
          </tr>
          <tr>
            <td>Family size</td>
            <td><label><img src="graphics/large family.jpg"/><br/>
            <input type="radio" name="r3">Large</input></label></td>
            <td>&nbsp;</td>
            <td><label><img src="graphics/small family.jpg"/><br/>
            <input type="radio" name="r3">Small</input></label></td>
          </tr>
          <tr>
            <td>Education</td>
            <td><label><img src="graphics/college.jpg"/><br/>
            <input type="radio" name="r4">College</input></label></td>
            <td><label><img src="graphics/some college.jpg"/><br/>
            <input type="radio" name="r4">Some College</input></label></td>
            <td><label><img src="graphics/no college.jpg"/><br/>
            <input type="radio" name="r4">No College</input></label></td>
          </tr>
          </tbody>
        </table>
        <p>
          <a href="dialog6.html" className="button">OK</a>
        </p>
      </div>
      <div className="frame">
        <table>
          <tbody>
          <tr>
            <td>&nbsp;</td>
            <td colspan="3" className="question">Second Record</td>
          </tr>
          <tr>
            <td>Fitness</td>
            <td><label><img src="graphics/excellent fitness.jpg"/><br/>
            <input type="radio" name="r1">Excellent</input></label></td>
            <td>&nbsp;</td>
            <td><label><img src="graphics/average fitness.jpg"/><br/>
            <input type="radio" name="r1">Average</input></label></td>
          </tr>
          <tr>
            <td>Parents health</td>
            <td><label><img src="graphics/excellent parents.jpg"/><br/>
            <input type="radio" name="r2">Excellent</input></label></td>
            <td>&nbsp;</td>
            <td><label><img src="graphics/fair parents.jpg"/><br/>
            <input type="radio" name="r2">Fair</input></label></td>
          </tr>
          <tr>
            <td>Family size</td>
            <td><label><img src="graphics/large family.jpg"/><br/>
            <input type="radio" name="r3">Large</input></label></td>
            <td>&nbsp;</td>
            <td><label><img src="graphics/small family.jpg"/><br/>
            <input type="radio" name="r3">Small</input></label></td>
          </tr>
          <tr>
            <td>Education</td>
            <td><label><img src="graphics/college.jpg"/><br/>
            <input type="radio" name="r4">College</input></label></td>
            <td><label><img src="graphics/some college.jpg"/><br/>
            <input type="radio" name="r4">Some College</input></label></td>
            <td><label><img src="graphics/no college.jpg"/><br/>
            <input type="radio" name="r4">No College</input></label></td>
          </tr>
          </tbody>
        </table>
        <p>
          <a href="dialog6.html" className="button">OK</a>
        </p>
      </div></div>;
  }
});

var OneRecordPerformance = React.createClass({

  // getInitialState: function() {
  //   return {mode: 0};
  // },

  render: function() {
      var state = this.state;
      var user = this.props.user;
      var app = this.props.app;
      var prompt = user.getPrompt();

  		return  <div className="frame">
        <table className="record">
          <tbody>
            <tr>
              <td colSpan="3" className="robot">Record #18 <b>Daisy Smith</b></td>
              <td className="robot">Gender: F</td>
            </tr>
          </tbody>
        </table>
        <table className="xxx">
          <tbody>
          <tr>
            <td>Fitness</td>
            <td><label className="dimmed"><img src="/img/excellent fitness.jpg"/><br/>
            Excellent</label></td>
            <td>&nbsp;</td>
            <td><label><img src="/img/average fitness.jpg"/><br/>
            Average</label></td>
          </tr>
          <tr>
            <td>Parents health</td>
            <td><label><img src="/img/excellent parents.jpg"/><br/>
            Excellent</label></td>
            <td>&nbsp;</td>
            <td><label className="dimmed"><img src="/img/fair parents.jpg"/><br/>
            Fair</label></td>
          </tr>
          <tr>
            <td>Family size</td>
            <td><label><img src="/img/large family.jpg"/><br/>
            Large</label></td>
            <td>&nbsp;</td>
            <td><label className="dimmed"><img src="/img/small family.jpg"/><br/>
            Small</label></td>
          </tr>
          <tr>
            <td>Education</td>
            <td><label className="dimmed"><img src="/img/college.jpg"/><br/>
            College</label></td>
            <td><label className="dimmed"><img src="/img/some college.jpg"/><br/>
            Some College</label></td>
            <td><label><img src="/img/no college.jpg"/><br/>
            No College</label></td>
          </tr>
          </tbody>
        </table>
        <p className="performance-level">Performance Level:
          <span className="grade">D</span>
        </p>
        </div>;
  }
});
